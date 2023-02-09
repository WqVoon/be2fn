package internal

import (
	"testing"
)

func TestLexer(t *testing.T) {
	expr := ` (a >= true) || (b <= "abc") && (c == false) && !(d != true)  `

	l := NewLexer(expr)
	if err := l.Parse(); err != nil {
		t.Fatalf("faild to parse expr, err: %v\n", err)
	}

	for _, param := range l.Params {
		t.Logf("%+v\n", param)
	}
}

func TestNotExpr(t *testing.T) {
	cases := []struct {
		Expr        string
		ShouldError bool
	}{
		{"!a", true},
		{"!!a", true},
		{"!!!a", true},
		{"!(a)", true},
		{"!!(a)", true},
		{"!!!(a)", true},
		{"!a && b", true},
		{"!(a==1 && b==1)", false},
		{"!!(a && b)", true},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}

func TestSubExpr(t *testing.T) {
	cases := []struct {
		Expr        string
		ShouldError bool
	}{
		{"-a", true},
		{"-1", false},
		{"-+1", true},
		{"--1", true},
		{"-(1)", true},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}

func TestBinaryExpr(t *testing.T) {
	cases := []struct {
		Expr        string
		ShouldError bool
	}{
		{"a == a", true},
		{"a == 1", false},
		{"1 == 1", true},
		{"1 == a", false},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}

func TestSlice(t *testing.T) {
	cases := []struct {
		Expr        string
		ShouldError bool
	}{
		{"[]int{}", false},
		{"[]int", false},
		{"[1]int{}", false},
		{"[]int{1}", false},
		{"[]int{1, 2}", false},
		{"[]int{1, 2, \"3\"}", true},
		{"[]uint{1}", true},
		{"[]string{}", false},
		{"[]string{\"1\"}", false},
		{"[]string{\"1\", 2}", true},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}

func TestFunc(t *testing.T) {
	cases := []struct {
		Expr        string
		ShouldError bool
	}{
		{"in()", true},
		{"in(1)", true},
		{"in(1, 2)", true},
		{"in(a, 2)", true},
		{"in(a, b)", true},
		{"in(a, []int)", true},
		{"in(a, []int{})", false},
		{"in(a, []string{})", false},
		{"in([]int{}, a)", true},
		{"in([]string{}, a)", true},
		{"no_such_func([]string{}, a)", true},
		{"no_such_func()", true},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}
