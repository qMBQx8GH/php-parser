package php84

import (
	"testing"

	"github.com/Demooon86/php-parser/internal/tester"
)

func TestConstantsInClasses(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP84()
	suite.Code = `<?php
		class Foo {
		    public const int CONSTANT = 1;
		    public const CONSTANT_2 = 2;
		}
    `

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Name: &ast.Identifier{
				Val: []byte("Foo"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtClassConstList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					ConstType: &ast.Name{
						Parts: []ast.Vertex{
							&ast.NamePart{
								Val: []byte("int"),
							},
						},
					},
					Consts: []ast.Vertex{
						&ast.StmtConstant{
							Name: &ast.Identifier{
								Val: []byte("CONSTANT"),
							},
							Expr: &ast.ScalarLnumber{
								Val: []byte("1"),
							},
						},
					},
				},
				&ast.StmtClassConstList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
					},
					Consts: []ast.Vertex{
						&ast.StmtConstant{
							Name: &ast.Identifier{
								Val: []byte("CONSTANT_2"),
							},
							Expr: &ast.ScalarLnumber{
								Val: []byte("2"),
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

func TestAsymmetricVisibility(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP84()
	suite.Code = `<?php
		class PhpVersion
		{
		    public private(set) string $version = '8.4';
		    private(set) int $count = 0;
		    protected(set) string $name;
		}
    `

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Name: &ast.Identifier{
				Val: []byte("PhpVersion"),
			},
			Stmts: []ast.Vertex{
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("public"),
						},
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
									Val: []byte("$version"),
								},
							},
							Expr: &ast.ScalarString{
								Val: []byte("'8.4'"),
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
								Val: []byte("int"),
							},
						},
					},
					Props: []ast.Vertex{
						&ast.StmtProperty{
							Var: &ast.ExprVariable{
								Name: &ast.Identifier{
									Val: []byte("$count"),
								},
							},
							Expr: &ast.ScalarLnumber{
								Val: []byte("0"),
							},
						},
					},
				},
				&ast.StmtPropertyList{
					Modifiers: []ast.Vertex{
						&ast.Identifier{
							Val: []byte("protected"),
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
									Val: []byte("$name"),
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
