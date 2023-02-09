package be2fn

import "github.com/wqvoon/be2fn/internal"

// 传给编译出的函数的参数，key 是字符串，value 是 interface{}，
// 目前 key 的类型在编译期根据二元表达式的另一个参数确定
type Kv = internal.Kv

// 将 expr 编译为一个可执行的函数，编译失败时返回错误原因
func Compile(expr string) (internal.Unit, error) {
	// 词法解析，生成组成逆波兰表达式的 token 序列
	lexer := internal.NewLexer(expr)
	if err := lexer.Parse(); err != nil {
		return nil, err
	}

	// 根据 lexer 解析出的 token 编译出最终的函数
	compiler := internal.NewCompiler(lexer)
	return compiler.Compile()
}
