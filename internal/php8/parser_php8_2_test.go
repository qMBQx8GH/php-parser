package php8_test

import (
	"testing"

	"github.com/VKCOM/php-parser/internal/tester"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/conf"
	"github.com/VKCOM/php-parser/pkg/errors"
	"github.com/VKCOM/php-parser/pkg/parser"
	"github.com/VKCOM/php-parser/pkg/version"
)

func TestClassReadonlyModifier(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP8()
	suite.Code = `<?php
readonly class Foo {
	public string $a;
}

final readonly class Boo {
	public string $a;
}
`

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Modifiers: []ast.Vertex{
				&ast.Identifier{
					Val: []byte("readonly"),
				},
			},
			Name: &ast.Identifier{
				Val: []byte("Foo"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					Type: &ast.Name{
						Parts: []ast.Vertex{
							&ast.NamePart{
								Val: []byte("string"),
							},
						},
					},
					Props: []ast.Vertex{
						&ast.StmtProperty{
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
						},
					},
				},
			},
		},
		&ast.StmtClass{
			Modifiers: []ast.Vertex{
				&ast.Identifier{
					Val: []byte("final"),
				},
				&ast.Identifier{
					Val: []byte("readonly"),
				},
			},
			Name: &ast.Identifier{
				Val: []byte("Boo"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					Type: &ast.Name{
						Parts: []ast.Vertex{
							&ast.NamePart{
								Val: []byte("string"),
							},
						},
					},
					Props: []ast.Vertex{
						&ast.StmtProperty{
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
						},
					},
				},
			},
		},
	},
},`

	suite.Run()
}

func TestDNFTypesAreParsed(t *testing.T) {
	codes := []string{
		`<?php function f((A&B)|null $x): (A&B)|null {}`,
		`<?php function f(A|(B&C)|D $x) {}`,
		`<?php function f((A&B)|(C&D) $x, (A&B)|string ...$rest): (A&B)|false {}`,
		`<?php function f((\Foo\A&B\C)|null $x) {}`,
		`<?php class X {
			public (A&B)|null $a;
			private static (A&B)|C $b;
			public function __construct(private (A&B)|null $c = null) {}
			public function m((A&B)|(C&D) $a): (A&B)|false {}
		}`,
		`<?php $f = fn((A&B)|null $x): (A&B)|null => $x;`,
		`<?php $f = function((A&B)|null $x) use ($y): (A&B)|null {};`,
		`<?php interface I { function f((A&B)|null $x): (A&B)|null; }`,
		// intersection types without parentheses still work
		`<?php function f(A&B $x): A&B {}`,
		// a reference parameter is not an intersection type
		`<?php function f(A &$x, (A&B)|null &...$rest) {}`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%s\nunexpected errors: %v", code, errs)
		}
	}
}

func TestDNFTypesAreRejected(t *testing.T) {
	codes := []string{
		// a parenthesized type must be an intersection of at least two types
		`<?php function f((A) $x) {}`,
		`<?php function f((A)|B $x) {}`,
		`<?php function f((A&B) $x) {}`,
		// nesting is not allowed
		`<?php function f(((A&B)|C)|D $x) {}`,
		`<?php function f((A|B)&C $x) {}`,
		`<?php function f((A&(B&C))|D $x) {}`,
		// the nullable prefix cannot be combined with DNF
		`<?php function f(?(A&B)|C $x) {}`,
		`<?php class X { public (A&B) $a; }`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) == 0 {
			t.Errorf("%s\nexpected a syntax error", code)
		}
	}
}

func TestDNFTypeAST(t *testing.T) {
	root, errs := parseWithErrors(`<?php function f((A&B)|null $x) {}`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	fn := root.(*ast.Root).Stmts[0].(*ast.StmtFunction)
	param := fn.Params[0].(*ast.Parameter)

	union, ok := param.Type.(*ast.Union)
	if !ok {
		t.Fatalf("param type is %T, want *ast.Union", param.Type)
	}
	if len(union.Types) != 2 {
		t.Fatalf("union has %d types, want 2", len(union.Types))
	}

	inter, ok := union.Types[0].(*ast.Intersection)
	if !ok {
		t.Fatalf("first element is %T, want *ast.Intersection", union.Types[0])
	}
	if inter.OpenParenthesisTkn == nil || inter.CloseParenthesisTkn == nil {
		t.Errorf("the parentheses of the intersection type are lost")
	}
	if len(inter.Types) != 2 {
		t.Errorf("intersection has %d types, want 2", len(inter.Types))
	}
	if got, want := inter.Position.StartPos, 17; got != want {
		t.Errorf("intersection starts at %d, want %d (the opening parenthesis)", got, want)
	}
	if got, want := inter.Position.EndPos, 22; got != want {
		t.Errorf("intersection ends at %d, want %d (the closing parenthesis)", got, want)
	}

	if _, ok := union.Types[1].(*ast.Name); !ok {
		t.Errorf("second element is %T, want *ast.Name", union.Types[1])
	}

	// Intersection types without parentheses have no parenthesis tokens.
	root, _ = parseWithErrors(`<?php function f(A&B $x) {}`)
	plain := root.(*ast.Root).Stmts[0].(*ast.StmtFunction).Params[0].(*ast.Parameter).Type.(*ast.Intersection)
	if plain.OpenParenthesisTkn != nil || plain.CloseParenthesisTkn != nil {
		t.Errorf("an intersection type without parentheses must not have parenthesis tokens")
	}
}

func TestParserAcceptsPHP82(t *testing.T) {
	for _, v := range []*version.Version{{Major: 8, Minor: 0}, {Major: 8, Minor: 1}, {Major: 8, Minor: 2}} {
		var errs []*errors.Error
		_, err := parser.Parse([]byte(`<?php function f((A&B)|null $x) {}`), conf.Config{
			Version:          v,
			ErrorHandlerFunc: func(e *errors.Error) { errs = append(errs, e) },
		})
		if err != nil {
			t.Errorf("version %d.%d: %v", v.Major, v.Minor, err)
		}
		if len(errs) != 0 {
			t.Errorf("version %d.%d: unexpected errors: %v", v.Major, v.Minor, errs)
		}
	}
}
