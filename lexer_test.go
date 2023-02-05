package be2fn

import (
	"testing"
)

func TestLexer(t *testing.T) {
	expr := ` (a >= true) || (b <= "abc") && (c == false) && !(d != true)  `

	l := NewLexer(expr)
	if err := l.Parse(); err != nil {
		t.Fatalf("faild to parse expr, err: %v\n", err)
	}

	for _, token := range l.Tokens {
		t.Logf("%+v\n", token)
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
