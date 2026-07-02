package php84

import (
	"testing"

	"github.com/VKCOM/php-parser/internal/tester"
)

func TestPropertyHooksGetSet(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP84()
	suite.Code = `<?php
		class Example
		{
		    public string $foo = 'default value' {
		        get {
		            return $this->foo;
		        }

		        set(string $value) {
		            $this->foo = strtolower($value);
		        }
		    }
		}
    `

	suite.Expected = `&ast.Root{
	Stmts: []ast.Vertex{
		&ast.StmtClass{
			Name: &ast.Identifier{
				Val: []byte("Example"),
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
									Val: []byte("$foo"),
								},
							},
							Expr: &ast.ScalarString{
								Val: []byte("'default value'"),
							},
							Hooks: []ast.Vertex{
								&ast.StmtPropertyHook{
									Stmts: []ast.Vertex{
										&ast.StmtReturn{
											Expr: &ast.ExprPropertyFetch{
												Var: &ast.ExprVariable{
													Name: &ast.Identifier{
														Val: []byte("$this"),
													},
												},
												Prop: &ast.Identifier{
													Val: []byte("foo"),
												},
											},
										},
									},
								},
								&ast.StmtPropertyHook{
									Params: []ast.Vertex{
										&ast.Parameter{
											Type: &ast.Name{
												Parts: []ast.Vertex{
													&ast.NamePart{
														Val: []byte("string"),
													},
												},
											},
											Var: &ast.ExprVariable{
												Name: &ast.Identifier{
													Val: []byte("$value"),
												},
											},
										},
									},
									Stmts: []ast.Vertex{
										&ast.StmtExpression{
											Expr: &ast.ExprAssign{
												Var: &ast.ExprPropertyFetch{
													Var: &ast.ExprVariable{
														Name: &ast.Identifier{
															Val: []byte("$this"),
														},
													},
													Prop: &ast.Identifier{
														Val: []byte("foo"),
													},
												},
												Expr: &ast.ExprFunctionCall{
													Function: &ast.Name{
														Parts: []ast.Vertex{
															&ast.NamePart{
																Val: []byte("strtolower"),
															},
														},
													},
													Args: []ast.Vertex{
														&ast.Argument{
															Expr: &ast.ExprVariable{
																Name: &ast.Identifier{
																	Val: []byte("$value"),
																},
															},
														},
													},
												},
											},
										},
									},
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

func TestPropertyHooksGetExprOnly(t *testing.T) {
	suite := tester.NewParserDumpTestSuite(t)
	suite.UsePHP84()
	suite.Code = `<?php
		class Foo
		{
		    public string $bar {
		        get => $this->bar . ' modified';
		    }
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
									Val: []byte("$bar"),
								},
							},
							Hooks: []ast.Vertex{
								&ast.StmtPropertyHook{
									Expr: &ast.ExprBinaryConcat{
										Left: &ast.ExprPropertyFetch{
											Var: &ast.ExprVariable{
												Name: &ast.Identifier{
													Val: []byte("$this"),
												},
											},
											Prop: &ast.Identifier{
												Val: []byte("bar"),
											},
										},
										Right: &ast.ScalarString{
											Val: []byte("' modified'"),
										},
									},
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
