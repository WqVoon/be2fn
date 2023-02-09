package internal

import "go/token"

type Token int

// 不使用原生的 golang token，方便做扩展，比如直接支持切片、布尔值等
const (
	INVALID Token = iota // 无效 token
	EOF                  // 文件尾

	// 可以被比较的东西
	IDENT     // 标识符，也就是通过 Unit 传进来的变量名
	INT       // 整数
	STRING    // 字符串
	BOOLEAN   // 布尔值
	INT_SLICE // 数字切片
	STR_SLICE // 字符串切片

	// 一元表达式操作符
	NOT // 非操作
	SUB // 减号，目前只用来处理负数

	// 二元表达式操作符
	LAND // &&
	LOR  // ||
	EQL  // ==
	NEQ  // !=
	LSS  // <
	LEQ  // <=
	GTR  // >
	GEQ  // >=
	FUNC // 函数调用
)

// 将 token 转换成对应的字符串表示
var token2String = map[Token]string{
	INVALID: "<invalid>",

	// 可以被比较的东西
	IDENT:     "ident",
	INT:       "int",
	STRING:    "string",
	BOOLEAN:   "boolean",
	INT_SLICE: "[]int",
	STR_SLICE: "[]string",

	// 一元表达式操作符
	NOT: "!",
	SUB: "-",

	// 二元表达式操作符
	LAND: "&&",
	LOR:  "||",
	EQL:  "==",
	NEQ:  "!=",
	LSS:  "<",
	LEQ:  "<=",
	GTR:  ">",
	GEQ:  ">=",
	FUNC: "func",
}

func (t Token) String() string {
	val, ok := token2String[t]
	if !ok {
		return "<invalid>"
	}
	return val
}

// 将 golang 中的 token 映射成 be2fn 的 token，
// 当且仅当 golang 的原生 token 在这个 map 中时，才是有效的 token
var golangToken2Token = map[token.Token]Token{
	token.EOF: EOF, // 文件尾

	// 可以被比较的东西
	token.IDENT:  IDENT,  // 标识符
	token.INT:    INT,    // 整数
	token.STRING: STRING, // 字符串

	// 一元表达式操作符
	token.NOT: NOT, // !
	token.SUB: SUB, // -

	// 二元表达式操作符
	token.LAND: LAND, // &&
	token.LOR:  LOR,  // ||
	token.EQL:  EQL,  // ==
	token.NEQ:  NEQ,  // !=
	token.LSS:  LSS,  // <
	token.LEQ:  LEQ,  // <=
	token.GTR:  GTR,  // >
	token.GEQ:  GEQ,  // >=
}
