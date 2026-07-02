package php83

import (
	"github.com/VKCOM/php-parser/internal/tester"
	"testing"
)

func TestConstantsInClasses(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP83()
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
