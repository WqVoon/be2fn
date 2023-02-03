package be2fn

import (
	"testing"
)

func TestLexer(t *testing.T) {
	expr := ` (a >= true) || (b <= "abc") && (c == false) && !d  `

	l := NewLexer(expr, 10)
	if err := l.Parse(); err != nil {
		t.Fatalf("faild to parse expr, err: %v\n", err)
	}

	for _, token := range l.Tokens {
		t.Logf("%+v\n", token)
	}
}
