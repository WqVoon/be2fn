package be2fn

import (
	"testing"
)

func TestLexer(t *testing.T) {
	expr := ` (a >= true) || (b <= "abc") && (c == false) && !(d != true)  `

	l := NewLexer(expr, 10)
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
		{"!(a && b)", false},
		{"!!(a && b)", true},
	}

	for i, c := range cases {
		hasError := (NewLexer(c.Expr, 10).Parse() != nil)

		if c.ShouldError && !hasError || !c.ShouldError && hasError {
			t.Fatalf("failed to test %d, expr: %q, shoudError: %v", i, c.Expr, c.ShouldError)
		}
	}
}
