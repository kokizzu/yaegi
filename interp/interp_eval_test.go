package interp_test

import (
	"reflect"
	"testing"

	"github.com/containous/dyngo/interp"
)

func TestEval0(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `var I int = 2`)

	t1 := evalCheck(t, i, `I`)
	if t1.Interface().(int) != 2 {
		t.Fatalf("expected 2, got %v", t1)
	}
}

func TestEval1(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `func Hello() string { return "hello" }`)

	v := evalCheck(t, i, `Hello`)

	f, ok := v.Interface().(func() string)
	if !ok {
		t.Fatal("conversion failed")
	}

	if s := f(); s != "hello" {
		t.Fatalf("expected hello, got %v", s)
	}
}

func TestEval2(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `package foo; var I int = 2`)

	t1 := evalCheck(t, i, `foo.I`)
	if t1.Interface().(int) != 2 {
		t.Fatalf("expected 2, got %v", t1)
	}
}

func TestEval3(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `package foo; func Hello() string { return "hello" }`)

	v := evalCheck(t, i, `foo.Hello`)
	f, ok := v.Interface().(func() string)
	if !ok {
		t.Fatal("conversion failed")
	}
	if s := f(); s != "hello" {
		t.Fatalf("expected hello, got %v", s)
	}
}

func TestEvalNil0(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `func getNil() error { return nil }`)

	v := evalCheck(t, i, `getNil()`)
	if !v.IsNil() {
		t.Fatalf("expected nil, got %v", v)
	}
}

func TestEvalNil1(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `
package bar

func New() func(string) error {
	return func(v string) error {
		return nil
	}
}
`)

	v := evalCheck(t, i, `bar.New()`)
	fn, ok := v.Interface().(func(string) error)
	if !ok {
		t.Fatal("conversion failed")
	}

	if res := fn("hello"); res != nil {
		t.Fatalf("expected nil, got %v", res)
	}
}

func TestEvalNil2(t *testing.T) {
	i := interp.New(interp.Opt{})
	_, err := i.Eval(`a := nil`)
	if err.Error() != "1:27: use of untyped nil" {
		t.Fatal("should have failed")
	}
}

func TestEvalComposite0(t *testing.T) {
	i := interp.New(interp.Opt{})
	evalCheck(t, i, `
type T struct {
	a, b, c, d, e, f, g, h, i, j, k, l, m, n string
	o map[string]int
	p []string
}

var a = T{
	o: map[string]int{"truc": 1, "machin": 2},
	p: []string{"hello", "world"},
}
`)
	v := evalCheck(t, i, `a.p[1]`)
	if v.Interface().(string) != "world" {
		t.Fatalf("expected world, got %v", v)
	}
}

func evalCheck(t *testing.T, i *interp.Interpreter, src string) reflect.Value {
	res, err := i.Eval(src)
	if err != nil {
		t.Fatal(err)
	}
	return res
}