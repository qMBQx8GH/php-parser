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
