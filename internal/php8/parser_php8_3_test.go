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

func TestDynamicClassConstantFetch(t *testing.T) {
	codes := []string{
		`<?php echo Config::{$name};`,
		`<?php echo Config::{"VER" . "SION"};`,
		`<?php echo $obj::{$name};`,
		`<?php echo static::{$name}, self::{$name}, parent::{$name}, \Foo\Bar::{$name};`,
		`<?php echo $a::{$b}::{$c};`,
		`<?php $x = Config::{$name} . Config::{$name2};`,
		`<?php class A { const X = 1; public function f($n) { return static::{$n}; } }`,
		// the old syntax still works
		`<?php echo Config::VERSION, Config::class, A::{$n}(), A::{$n}(1, 2), A::$n, A::$$n;`,
		`<?php echo A::{$n}()::{$m};`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%s\nunexpected errors: %v", code, errs)
		}
	}
}

func TestDynamicClassConstantFetchIsRejected(t *testing.T) {
	codes := []string{
		`<?php echo Config::{};`,
		`<?php echo Config::{$a, $b};`,
		`<?php echo Config::{$a;`,
		`<?php Config::{$name} = 1;`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) == 0 {
			t.Errorf("%s\nexpected a syntax error", code)
		}
	}
}

func TestDynamicClassConstantFetchAST(t *testing.T) {
	expr := firstExpr(t, `<?php $x = Config::{$name};`)

	fetch, ok := expr.(*ast.ExprClassConstFetch)
	if !ok {
		t.Fatalf("got %T, want *ast.ExprClassConstFetch", expr)
	}
	brackets, ok := fetch.Const.(*ast.ExprBrackets)
	if !ok {
		t.Fatalf("Const is %T, want *ast.ExprBrackets", fetch.Const)
	}
	if v, ok := brackets.Expr.(*ast.ExprVariable); !ok || v.Name.(*ast.Identifier).Value[0] != '$' {
		t.Errorf("the expression in the braces is %T", brackets.Expr)
	}

	// The static constant has the identifier.
	expr = firstExpr(t, `<?php $x = Config::VERSION;`)
	if _, ok := expr.(*ast.ExprClassConstFetch).Const.(*ast.Identifier); !ok {
		t.Errorf("a static constant must have the *ast.Identifier")
	}

	// A call still is a call.
	expr = firstExpr(t, `<?php $x = Config::{$name}();`)
	if _, ok := expr.(*ast.ExprStaticCall); !ok {
		t.Errorf("got %T, want *ast.ExprStaticCall", expr)
	}
}

func TestAnonymousReadonlyClass(t *testing.T) {
	codes := []string{
		`<?php $o = new readonly class { public int $x; };`,
		`<?php $o = new readonly class(5) { public function __construct(public int $x) {} };`,
		`<?php $o = new readonly class(5) extends Base implements I, J { };`,
		`<?php $o = new #[Attr] readonly class { };`,
		`<?php $o = new class { };`,
		`<?php $o = new class(1, 2) extends Base { };`,
		`<?php f(new readonly class {});`,
	}

	for _, code := range codes {
		if _, errs := parseWithErrors(code); len(errs) != 0 {
			t.Errorf("%s\nunexpected errors: %v", code, errs)
		}
	}

	for _, code := range []string{
		`<?php $o = new readonly {};`,
		`<?php $o = new readonly readonly class {};`,
		`<?php $o = new final class {};`,
	} {
		if _, errs := parseWithErrors(code); len(errs) == 0 {
			t.Errorf("%s\nexpected a syntax error", code)
		}
	}
}

func TestAnonymousReadonlyClassAST(t *testing.T) {
	expr := firstExpr(t, `<?php $o = new readonly class(5) { public int $x; };`)

	class := expr.(*ast.ExprNew).Class.(*ast.StmtClass)
	if class.Name != nil {
		t.Errorf("an anonymous class has the name")
	}
	if len(class.Modifiers) != 1 || string(class.Modifiers[0].(*ast.Identifier).Value) != "readonly" {
		t.Errorf("the readonly modifier is lost: %v", class.Modifiers)
	}
	if len(class.Args) != 1 {
		t.Errorf("the arguments are lost: %v", class.Args)
	}

	plain := firstExpr(t, `<?php $o = new class { };`).(*ast.ExprNew).Class.(*ast.StmtClass)
	if len(plain.Modifiers) != 0 {
		t.Errorf("the plain anonymous class has the modifiers")
	}
}
