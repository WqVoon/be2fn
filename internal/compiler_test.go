package internal

import "testing"

func TestNotAndOr(t *testing.T) {
	lex := NewLexer(`
		(a > 0 && a < 10) ||
		(a > -10 && a < 0) ||
		b == true ||
		!(c == "") ||
		(d <= 1 && d >= 1) ||
		(e != 0)
	`)

	if err := lex.Parse(); err != nil {
		t.Fatal("faild to call Parse, err:", err)
	}

	for _, tkn := range lex.Tokens {
		t.Log(tkn)
	}

	fn, err := NewCompiler(lex).Compile()
	if err != nil {
		t.Fatal("failed to call Compile, err:", err)
	}

	cases := []struct {
		arg Kv
		ret bool
	}{
		// 测试 a
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 0}, ret: false},
		{arg: Kv{"a": 1, "b": false, "c": "", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": 2, "b": false, "c": "", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": 10, "b": false, "c": "", "d": 0, "e": 0}, ret: false},
		{arg: Kv{"a": -1, "b": false, "c": "", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": -2, "b": false, "c": "", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": -10, "b": false, "c": "", "d": 0, "e": 0}, ret: false},

		// 测试 b
		{arg: Kv{"a": 0, "b": true, "c": "", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 0}, ret: false},

		// 测试 c
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 0}, ret: false},
		{arg: Kv{"a": 0, "b": false, "c": "123", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": 0, "b": false, "c": "456", "d": 0, "e": 0}, ret: true},
		{arg: Kv{"a": 0, "b": false, "c": "789", "d": 0, "e": 0}, ret: true},

		// 测试 d
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 0}, ret: false},
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 1, "e": 0}, ret: true},

		// 测试 e
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 0}, ret: false},
		{arg: Kv{"a": 0, "b": false, "c": "", "d": 0, "e": 1}, ret: true},
	}

	for idx, c := range cases {
		fnRet, err := fn(c.arg)
		if err != nil {
			t.Fatalf("failed to call fn, err: %v", err)
		}

		if fnRet != c.ret {
			t.Fatalf("failed to pass case(%d), want(%v), got(%v)", idx, c.ret, fnRet)
		}
	}
}
