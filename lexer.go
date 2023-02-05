package be2fn

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// 表达式中允许出现的 token 类型
var isValidTokenType = map[token.Token]bool{
	token.EOF: true, // 文件尾

	// 可以被比较的东西
	token.IDENT:  true, // 标识符
	token.INT:    true, // 整数
	token.STRING: true, // 字符串

	// 一元表达式操作符
	token.NOT: true, // !
	token.SUB: true, // -，目前只用来处理负数，类似 `-1` 这种后面直接跟数字的

	// 二元表达式操作符
	token.LAND: true, // &&
	token.LOR:  true, // ||
	token.EQL:  true, // ==
	token.NEQ:  true, // !=
	token.LSS:  true, // <
	token.LEQ:  true, // <=
	token.GTR:  true, // >
	token.GEQ:  true, // >=

}

type Token struct {
	Typ token.Token
	Val string

	IsBoolean bool // 当前 token 是否是一个布尔值，如果是则 Typ 字段为 token.ILLEGAL
	BoolVal   bool // IsBoolean 为 true 时，BoolVal 代表真正的布尔值内容

	IntVal int // 当前 token 是数字时，这里保存实际的值
}

type Lexer struct {
	Err        error   // 解析时遇到的错误
	HasParsed  bool    // 是否已经解析过
	SourceCode string  // 原表达式
	Tokens     []Token // 解析的结果，是一个合法的逆波兰表达式的 token 序列

	ExecWhenWalk func(node ast.Node) // 可以自定义的函数，针对 AST 上的每个节点都会执行
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{
		SourceCode: sourceCode,
		Tokens:     make([]Token, 0, len(sourceCode)/2),
	}
}

func NewLexerWithTokenSize(sourceCode string, tokenSize int) *Lexer {
	return &Lexer{
		SourceCode: sourceCode,
		Tokens:     make([]Token, 0, tokenSize),
	}
}

// 解析 SourceCode 为 Tokens
func (l *Lexer) Parse() error {
	if l.HasParsed {
		return l.Err
	}
	defer func() { l.HasParsed = true }()

	expr, err := parser.ParseExpr(l.SourceCode)
	if err != nil {
		return err
	}

	return l.walk(expr)
}

// 后序遍历 AST
func (l *Lexer) walk(node ast.Node) error {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if err := l.walk(n.X); err != nil {
			return err
		}
		if err := l.walk(n.Y); err != nil {
			return err
		}

	case *ast.ParenExpr:
		l.walk(n.X)
		return l.Err

	case *ast.UnaryExpr:
		if err := l.walk(n.X); err != nil {
			return err
		}
	}

	l.handleOneNode(node)
	return l.Err
}

// 处理一个 ast 节点
func (l *Lexer) handleOneNode(node ast.Node) (shouldStop bool) {
	if l.ExecWhenWalk != nil { // 执行自定义函数
		l.ExecWhenWalk(node)
	}

	switch expr := node.(type) {
	case *ast.BinaryExpr:
		// fmt.Println(expr.Op)
		return !l.handleBinaryExpr(expr)

	case *ast.UnaryExpr:
		// fmt.Println(expr.Op)
		return !l.handleUnaryExpr(expr)

	case *ast.BasicLit:
		// fmt.Println(expr.Value)
		return !l.handleBasicLit(expr)

	case *ast.Ident:
		// fmt.Println(expr.Name)
		return !l.handleIdent(expr)
	}

	return
}

// 处理二元表达式
func (l *Lexer) handleBinaryExpr(be *ast.BinaryExpr) (isValid bool) {
	if !isValidTokenType[be.Op] {
		l.Err = invalidTokenError(be.Op, be.OpPos)
		return false
	}

	if be.Op == token.LAND || be.Op == token.LOR { // and/or 的子表达式必须为二元表达式或括号包裹的二元表达式
		_, xIsBE := be.X.(*ast.BinaryExpr)
		if !xIsBE && !isParenExprwithBinaryExpr(be.X) && !isUnaryExprWithNotOp(be.X) {
			l.Err = fmt.Errorf("`%s`'s subExpr must be BinaryExpr or ParenExpr(with BinaryExpr) or UnaryExpr(with `not` op), err at %v", be.Op, be.OpPos)
			return false
		}

		_, yIsBE := be.Y.(*ast.BinaryExpr)
		if !yIsBE && !isParenExprwithBinaryExpr(be.Y) && !isUnaryExprWithNotOp(be.Y) {
			l.Err = fmt.Errorf("`%s`'s subExpr must be BinaryExpr or ParenExpr(with BinaryExpr) or UnaryExpr(with `not` op), err at %v", be.Op, be.OpPos)
			return false
		}
	}

	if isBasicLit(be.X) && isBasicLit(be.Y) { // 不支持操作数均为常量的判断
		l.Err = fmt.Errorf("both subExpr of `%s` is BasicLit, err at %v", be.Op, be.OpPos)
		return false
	}

	if isIdent(be.X) && isIdent(be.Y) { // 不支持操作数均为变量或布尔值的判断
		xIsBool, yIsBool := isBoolIdent(be.X.(*ast.Ident)), isBoolIdent(be.Y.(*ast.Ident))
		if (xIsBool && yIsBool) || (!xIsBool && !yIsBool) {
			l.Err = fmt.Errorf("both subExpr of `%s` is BoolValue or Ident, err at %v", be.Op, be.OpPos)
			return false
		}
	}

	l.Tokens = append(l.Tokens, Token{Typ: be.Op, Val: be.Op.String()})
	return true
}

// 处理一元表达式
func (l *Lexer) handleUnaryExpr(ue *ast.UnaryExpr) (isValid bool) {
	if !isValidTokenType[ue.Op] {
		l.Err = invalidTokenError(ue.Op, ue.OpPos)
		return false
	}

	switch ue.Op {
	case token.NOT:
		if !isParenExprwithBinaryExpr(ue.X) { // not 的子表达式必须是括号包裹的二元表达式
			l.Err = fmt.Errorf("`not`'s subExpr must be ParenExpr(with BinaryExpr), err at %v", ue.OpPos)
			return false
		}

	case token.SUB:
		basicLit, ok := ue.X.(*ast.BasicLit)
		if !ok || basicLit.Kind != token.INT { // 负号后面必须跟着一个数字常量
			l.Err = fmt.Errorf("`-`'s subExpr must be number, err at %v", ue.OpPos)
		}
	}

	l.Tokens = append(l.Tokens, Token{Typ: ue.Op, Val: ue.Op.String()})
	return true
}

// 处理数字或字符串
func (l *Lexer) handleBasicLit(lt *ast.BasicLit) (isValid bool) {
	if !isValidTokenType[lt.Kind] {
		l.Err = invalidTokenError(lt.Kind, lt.Pos())
		return false
	}

	switch lt.Kind {
	case token.INT:
		intVal, _ := strconv.ParseInt(lt.Value, 10, 64)
		l.Tokens = append(l.Tokens, Token{Typ: lt.Kind, Val: lt.Value, IntVal: int(intVal)})

	case token.STRING:
		l.Tokens = append(l.Tokens, Token{Typ: lt.Kind, Val: strings.Trim(lt.Value, `"`)})

	default:
		return false
	}

	return true
}

// 处理标识符或布尔值
func (l *Lexer) handleIdent(it *ast.Ident) (isValid bool) {
	if isBoolIdent(it) {
		boolVal, _ := strconv.ParseBool(it.Name)
		l.Tokens = append(l.Tokens, Token{Typ: token.IDENT, IsBoolean: true, BoolVal: boolVal})
	} else {
		l.Tokens = append(l.Tokens, Token{Typ: token.IDENT, Val: it.Name})
	}

	return true
}

// 生成无效 token 的错误信息
func invalidTokenError(t token.Token, pos token.Pos) error {
	return fmt.Errorf("invalid token(%q) at position(%v)", t, pos)
}

// 判断 expr 是否是标识符
func isIdent(expr ast.Expr) bool {
	_, isIdent := expr.(*ast.Ident)
	return isIdent
}

// 判断 expr 是否是字符串或数字
func isBasicLit(expr ast.Expr) bool {
	_, isBasicLit := expr.(*ast.BasicLit)
	return isBasicLit
}

func isBoolIdent(ident *ast.Ident) bool {
	return ident.Name == "true" || ident.Name == "false"
}

// 判断 expr 是否为包含 BinaryExpr 的 ParenExpr，给 not/and/or 用，
// 因为目前 not 操作符只能处理这种表达式，and/or 只能处理这种以及 BinaryExpr
func isParenExprwithBinaryExpr(expr ast.Expr) bool {
	paranExpr, ok := expr.(*ast.ParenExpr)
	if !ok {
		return false
	}

	_, ok = paranExpr.X.(*ast.BinaryExpr)
	return ok
}

// 判断 expr 是否为 not 的一元表达式
func isUnaryExprWithNotOp(expr ast.Expr) bool {
	be, ok := expr.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	return be.Op == token.NOT
}
