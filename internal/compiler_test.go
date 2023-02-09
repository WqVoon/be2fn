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

	for _, tkn := range lex.Params {
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

func TestIn(t *testing.T) {
	type Pair struct {
		Vars Kv
		Ret  bool
	}

	cases := []struct {
		Expr     string
		SubCases []Pair
	}{
		// in
		{"in(a, []int{1,2,3})", []Pair{
			{Kv{"a": 0}, false},
			{Kv{"a": 1}, true},
			{Kv{"a": 2}, true},
			{Kv{"a": 3}, true},
			{Kv{"a": 4}, false},
		}},

		{`in(a, []string{"1","2","3"})`, []Pair{
			{Kv{"a": "0"}, false},
			{Kv{"a": "1"}, true},
			{Kv{"a": "2"}, true},
			{Kv{"a": "3"}, true},
			{Kv{"a": "4"}, false},
		}},

		// not in
		{"! in(a, []int{1,2,3})", []Pair{
			{Kv{"a": 0}, true},
			{Kv{"a": 1}, false},
			{Kv{"a": 2}, false},
			{Kv{"a": 3}, false},
			{Kv{"a": 4}, true},
		}},

		{`! in(a, []string{"1","2","3"})`, []Pair{
			{Kv{"a": "0"}, true},
			{Kv{"a": "1"}, false},
			{Kv{"a": "2"}, false},
			{Kv{"a": "3"}, false},
			{Kv{"a": "4"}, true},
		}},
	}

	for idx, c := range cases {
		t.Logf("start to test case %d", idx)

		lex := NewLexer(c.Expr)
		if err := lex.Parse(); err != nil {
			t.Fatalf("faild to parse %q, err: %v", c.Expr, err)
		}

		for _, tkn := range lex.Params {
			t.Log(tkn)
		}

		fn, err := NewCompiler(lex).Compile()
		if err != nil {
			t.Fatalf("failed to compile %q, err: %v", c.Expr, err)
		}

		for _, pair := range c.SubCases {
			ret, err := fn(pair.Vars)
			if err != nil {
				t.Fatalf("failed to call fn for expr(%v) with kv(%v), err: %v", c.Expr, pair.Vars, err)
			}

			if ret != pair.Ret {
				t.Fatalf("failed to call fn for expr(%v) with kv(%v), shouldRet: %v", c.Expr, pair.Vars, pair.Ret)
			}
		}

		t.Logf("test case %d pass", idx)
	}
}
