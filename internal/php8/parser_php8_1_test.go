package php8_test

import (
	"strings"
	"testing"

	"github.com/VKCOM/php-parser/internal/tester"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/conf"
	"github.com/VKCOM/php-parser/pkg/errors"
	"github.com/VKCOM/php-parser/pkg/parser"
	"github.com/VKCOM/php-parser/pkg/version"
)

func TestReadonlyModifier(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP8()
	suite.Code = `<?php
class Foo {
	readonly string $a;
	private readonly string $a;
	private string $a;
	private readonly $a = 100;

	public function __construct(
		readonly string $a,
		private readonly string $a,
		private string $a,
		private readonly $a = 100,
	) {}
}
`

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Name: &ast.Identifier{
				Val: []byte("Foo"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("readonly"),
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
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("private"),
						},
						&ast.Identifier{
							Val: []byte("readonly"),
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
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("private"),
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
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("private"),
						},
						&ast.Identifier{
							Val: []byte("readonly"),
						},
					},
					Props: []ast.Vertex{
						&ast.StmtProperty{
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
							Expr: &ast.ScalarLnumber{
								Val: []byte("100"),
							},
						},
					},
				},
				&ast.StmtClassMethod{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					Name: &ast.Identifier{
						Val: []byte("__construct"),
					},
					Params: []ast.Vertex{
						&ast.Parameter{
							Modifiers: []ast.Vertex{
								&ast.Identifier{
									Val: []byte("readonly"),
								},
							},
							Type: &ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("string"),
									},
								},
							},
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
						},
						&ast.Parameter{
							Modifiers: []ast.Vertex{
								&ast.Identifier{
									Val: []byte("private"),
								},
								&ast.Identifier{
									Val: []byte("readonly"),
								},
							},
							Type: &ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("string"),
									},
								},
							},
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
						},
						&ast.Parameter{
							Modifiers: []ast.Vertex{
								&ast.Identifier{
									Val: []byte("private"),
								},
							},
							Type: &ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("string"),
									},
								},
							},
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
						},
						&ast.Parameter{
							Modifiers: []ast.Vertex{
								&ast.Identifier{
									Val: []byte("private"),
								},
								&ast.Identifier{
									Val: []byte("readonly"),
								},
							},
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$a"),
								},
							},
							DefaultValue: &ast.ScalarLnumber{
								Val: []byte("100"),
							},
						},
					},
					Stmt: &ast.StmtStmtList{
						Stmts: []ast.Vertex{},
					},
				},
			},
		},
	},
},`

	suite.Run()
}

func TestNeverType(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP8()
	suite.Code = `<?php
function f(): never {}
`

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtFunction{
			Name: &ast.Identifier{
				Val: []byte("f"),
			},
			ReturnType: &ast.Name{
				Parts: []ast.Vertex{
					&ast.NamePart{
						Val: []byte("never"),
					},
				},
			},
			Stmts: []ast.Vertex{},
		},
	},
},`

	suite.Run()
}

func TestEnum(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP8()
	suite.Code = `<?php
enum A {}
enum B implements Bar, Baz {
}
enum C: int implements Bar {}
`

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtEnum{
			Name: &ast.Identifier{
				Val: []byte("A"),
			},
			Stmts: []ast.Vertex{},
		},
		&ast.StmtEnum{
			Name: &ast.Identifier{
				Val: []byte("B"),
			},
			Implements: []ast.Vertex{
				&ast.Name{
					Parts: []ast.Vertex{
						&ast.NamePart{
							Val: []byte("Bar"),
						},
					},
				},
				&ast.Name{
					Parts: []ast.Vertex{
						&ast.NamePart{
							Val: []byte("Baz"),
						},
					},
				},
			},
			Stmts: []ast.Vertex{},
		},
		&ast.StmtEnum{
			Name: &ast.Identifier{
				Val: []byte("C"),
			},
			Type: &ast.Name{
				Parts: []ast.Vertex{
					&ast.NamePart{
						Val: []byte("int"),
					},
				},
			},
			Implements: []ast.Vertex{
				&ast.Name{
					Parts: []ast.Vertex{
						&ast.NamePart{
							Val: []byte("Bar"),
						},
					},
				},
			},
			Stmts: []ast.Vertex{},
		},
	},
},`

	suite.Run()
}

func TestIntersectionTypes(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP8()
	suite.Code = `<?php
class Test {
    public A&B $prop;
}

function test(A&B $a): A&B {}
`

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Name: &ast.Identifier{
				Val: []byte("Test"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					Type: &ast.Intersection{
						Types: []ast.Vertex{
							&ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("A"),
									},
								},
							},
							&ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("B"),
									},
								},
							},
						},
					},
					Props: []ast.Vertex{
						&ast.StmtProperty{
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$prop"),
								},
							},
						},
					},
				},
			},
		},
		&ast.StmtFunction{
			Name: &ast.Identifier{
				Val: []byte("test"),
			},
			Params: []ast.Vertex{
				&ast.Parameter{
					Type: &ast.Intersection{
						Types: []ast.Vertex{
							&ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("A"),
									},
								},
							},
							&ast.Name{
								Parts: []ast.Vertex{
									&ast.NamePart{
										Val: []byte("B"),
									},
								},
							},
						},
					},
					Var: &ast.ExprVariable{
						Name: &ast.Identifier{
							Val: []byte("$a"),
						},
					},
				},
			},
			ReturnType: &ast.Intersection{
				Types: []ast.Vertex{
					&ast.Name{
						Parts: []ast.Vertex{
							&ast.NamePart{
								Val: []byte("A"),
							},
						},
					},
					&ast.Name{
						Parts: []ast.Vertex{
							&ast.NamePart{
								Val: []byte("B"),
							},
						},
					},
				},
			},
			Stmts: []ast.Vertex{},
		},
	},
},`

	suite.Run()
}

// Reserved keywords are valid names of enum cases, like they are for class constants.
func TestEnumCaseReservedKeywordNames(t *testing.T) {
	names := []string{"Use", "Empty", "Default", "List", "Array", "Class", "Fn", "Match", "New", "Print", "Static", "Public", "Enum", "None"}

	for _, name := range names {
		src := "<?php\nenum E: int { case " + name + " = 1; }"

		root, errs := parseWithErrors(src)
		if len(errs) != 0 {
			t.Errorf("case %s: unexpected errors: %v", name, errs)
			continue
		}

		enum := root.(*ast.Root).Stmts[0].(*ast.StmtEnum)
		enumCase := enum.Stmts[0].(*ast.EnumCase)
		if got := string(enumCase.Name.(*ast.Identifier).Value); got != name {
			t.Errorf("case %s: got name %q", name, got)
		}
	}
}

// `enum` is a keyword only when it is followed by the name of an enum.
func TestEnumAsIdentifier(t *testing.T) {
	codes := []string{
		"<?php class Enum {}",
		"<?php class Enum extends Foo {}",
		"<?php class Enum implements Foo {}",
		"<?php $a = Enum::X;",
		"<?php $a = new Enum();",
		"<?php use A\\Enum\\B; namespace X\\Enum;",
		"<?php function enum() {} enum();",
		"<?php $x->enum = 1; $x->enum();",
		"<?php enum Foo {} class Bar extends Enum {}",
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", code, errs)
		}
	}
}

func parseWithErrors(src string) (ast.Vertex, []*errors.Error) {
	var errs []*errors.Error
	root, _ := parser.Parse([]byte(src), conf.Config{
		Version:          &version.Version{Major: 8, Minor: 1},
		ErrorHandlerFunc: func(e *errors.Error) { errs = append(errs, e) },
	})
	return root, errs
}

// A reference to a static property inside of an array literal: `&self::$x`.
func TestArrayItemReferenceToStaticProperty(t *testing.T) {
	codes := []string{
		"<?php $a = ['k' => &self::$c];",
		"<?php $a = [&self::$c];",
		"<?php $a = ['k' => &Foo::$c, &static::$d];",
		"<?php $a = ['k' => &$c, &$d];",
		"<?php $a = array('k' => &self::$c);",
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", code, errs)
		}
	}
}

// PHP accepts the `<?php` open tag at the very end of a file.
func TestOpenTagAtEndOfFile(t *testing.T) {
	codes := []string{
		"<?php\nfoo();\n?>x<?php",
		"<?php\nfoo();\n?>x<?PHP",
		"x<?php",
		"<?php",
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", code, errs)
		}
	}
}

func firstExpr(t *testing.T, code string) ast.Vertex {
	t.Helper()

	root, errs := parseWithErrors(code)
	if len(errs) != 0 {
		t.Fatalf("%q: unexpected errors: %v", code, errs)
	}
	stmt := root.(*ast.Root).Stmts[0].(*ast.StmtExpression)
	return stmt.Expr.(*ast.ExprAssign).Expr
}

func partKinds(parts []ast.Vertex) []string {
	var kinds []string
	for _, p := range parts {
		switch p := p.(type) {
		case *ast.ScalarEncapsedStringPart:
			kinds = append(kinds, "text:"+string(p.Value))
		case *ast.ExprVariable:
			kinds = append(kinds, "var")
		case *ast.ExprPropertyFetch:
			kinds = append(kinds, "prop")
		case *ast.ExprNullsafePropertyFetch:
			kinds = append(kinds, "nullsafe")
		case *ast.ExprArrayDimFetch:
			kinds = append(kinds, "dim")
		default:
			kinds = append(kinds, "other")
		}
	}
	return kinds
}

// In the simple interpolation syntax only the first property is a part of the expression.
func TestSimpleInterpolationTakesOnlyOneProperty(t *testing.T) {
	tests := map[string][]string{
		`<?php $s = "[U=$id]$o->name[/U]";`: {"text:[U=", "var", "text:]", "prop", "text:[/U]"},
		`<?php $s = "$o->name[0]";`:         {"prop", "text:[0]"},
		`<?php $s = "$o->a->b";`:            {"prop", "text:->b"},
		`<?php $s = "$o?->a[0]";`:           {"nullsafe", "text:[0]"},
		`<?php $s = "$a[0]->b";`:            {"dim", "text:->b"},
		`<?php $s = "x $o->a y";`:           {"text:x ", "prop", "text: y"},
	}

	for code, want := range tests {
		got := partKinds(firstExpr(t, code).(*ast.ScalarEncapsed).Parts)
		if strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("%s\n got  %v\n want %v", code, got, want)
		}
	}
}

// A quote or a dollar is escaped by an odd number of backslashes only.
func TestEscapedQuoteAndDollarParity(t *testing.T) {
	// "\\\"\"$e": text `\\\"\"` and the variable $e.
	got := partKinds(firstExpr(t, `<?php $s = "\\\"\"$e";`).(*ast.ScalarEncapsed).Parts)
	if want := `text:\\\"\"|var`; strings.Join(got, "|") != want {
		t.Errorf("got %v, want %s", got, want)
	}

	// "\\$e": an escaped backslash and the variable.
	got = partKinds(firstExpr(t, `<?php $s = "\\$e";`).(*ast.ScalarEncapsed).Parts)
	if want := `text:\\|var`; strings.Join(got, "|") != want {
		t.Errorf("got %v, want %s", got, want)
	}

	// "\\\$e": an escaped backslash and an escaped dollar, no variable at all.
	if _, ok := firstExpr(t, `<?php $s = "\\\$e";`).(*ast.ScalarString); !ok {
		t.Errorf(`"\\\$e" must be a string without interpolation`)
	}

	// The string from the Bitrix code base.
	code := `<?php $s = "\"(?:[^\\\"\\r\\\\]|\\\\.|$white_space)*\"$white_space*";`
	if _, ok := firstExpr(t, code).(*ast.ScalarEncapsed); !ok {
		t.Errorf("%s must be an interpolated string", code)
	}
}

// A lone \r is a line break (old Mac files).
func TestCarriageReturnLineBreaks(t *testing.T) {
	for _, code := range []string{
		"<?php\r$a = 1;\r$b = 2;",
		"<?php $a = 1;\r",
		"<?\r$a = 1;?>x\r<?\r$b = 2;",
		"<?php\r\r$a = 1; // comment\r$b = 2;",
	} {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", code, errs)
		}
	}

	root, _ := parseWithErrors("<?php\r$a = 1;\r$b = 2;")
	stmts := root.(*ast.Root).Stmts
	if got := stmts[0].GetPosition().StartLine; got != 2 {
		t.Errorf("first statement is at line %d, want 2", got)
	}
	if got := stmts[1].GetPosition().StartLine; got != 3 {
		t.Errorf("second statement is at line %d, want 3", got)
	}
}
