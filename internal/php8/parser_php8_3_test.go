package php8_test

import (
	"testing"

	"github.com/VKCOM/php-parser/pkg/ast"
)

func TestTypedClassConstantsAreParsed(t *testing.T) {
	codes := []string{
		`<?php class A { const int X = 1; }`,
		`<?php class A { public const string X = "a", Y = "b"; }`,
		`<?php class A { final protected const ?int X = null; private const int|string Z = 1; }`,
		`<?php class A { const array LIST = [1, 2]; const \Foo\Bar B = null; const self S = null; const static T = null; }`,
		`<?php interface I { public const bool FLAG = true; const mixed ANY = 1; }`,
		`<?php trait T { public const float PI = 3.14; }`,
		`<?php enum E: int { const int DEFAULT = 1; case A = 1; }`,
		`<?php class A { #[Foo] public const callable F = 'x'; }`,
		`<?php class A { const (A&B)|null C = null; }`,
		// the untyped constants still work, including the names of the types
		`<?php class A { const X = 1; const Y = 2, Z = 3; const int = 1; const string = 'a'; const array = []; }`,
		`<?php class A { const FOO = 1; const BAR = self::FOO; }`,
		// trait adaptations use the same names
		`<?php class A { use T { foo as bar; baz as protected; } }`,
		`<?php class A { use T, U { T::foo insteadof U; U::foo as qux; } }`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%s\nunexpected errors: %v", code, errs)
		}
	}
}

func TestTypedClassConstantsAreRejected(t *testing.T) {
	codes := []string{
		`<?php class A { const int = ; }`,
		`<?php class A { const int X; }`,
		`<?php class A { const int string X = 1; }`,
		`<?php class A { const int; }`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) == 0 {
			t.Errorf("%s\nexpected a syntax error", code)
		}
	}
}

func TestTypedClassConstantAST(t *testing.T) {
	root, errs := parseWithErrors(`<?php class A { public const ?int X = 1, Y = 2; const Z = 3; }`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	class := root.(*ast.Root).Stmts[0].(*ast.StmtClass)

	typed := class.Stmts[0].(*ast.StmtClassConstList)
	if _, ok := typed.Type.(*ast.Nullable); !ok {
		t.Fatalf("type is %T, want *ast.Nullable", typed.Type)
	}
	if len(typed.Consts) != 2 {
		t.Errorf("%d constants, want 2", len(typed.Consts))
	}

	untyped := class.Stmts[1].(*ast.StmtClassConstList)
	if untyped.Type != nil {
		t.Errorf("an untyped constant has the type %T", untyped.Type)
	}
}
