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
	token.ILLEGAL: true, // 目前用来代表布尔值
	token.EOF:     true, // 文件尾

	token.IDENT:  true, // 标识符
	token.INT:    true, // 整数
	token.STRING: true, // 字符串

	token.LAND: true, // &&
	token.LOR:  true, // ||
	token.NOT:  true, // !
	token.EQL:  true, // ==
	token.NEQ:  true, // !=
	token.LSS:  true, // <
	token.LEQ:  true, // <=
	token.GTR:  true, // >
	token.GEQ:  true, // >=

	token.SUB: true, // -，目前只用来处理负数，类似 `-1` 这种后面直接跟数字的
}

type Token struct {
	Typ token.Token
	Val string

	IsBoolean bool
	BoolVal   bool
}

type Lexer struct {
	Err        error
	HasParsed  bool
	SourceCode string
	Tokens     []Token

	ExecWhenWalk func(node ast.Node)
}

func NewLexer(sourceCode string, tokenSize int) *Lexer {
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

	l.walk(expr)
	return l.Err
}

// 后序遍历 AST
func (l *Lexer) walk(node ast.Node) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		l.walk(n.X)
		l.walk(n.Y)

	case *ast.ParenExpr:
		l.walk(n.X)
		return

	case *ast.UnaryExpr:
		l.walk(n.X)
	}

	if l.handleOneNode(node) {
		return
	}
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

	if be.Op == token.LAND || be.Op == token.LOR {
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
		if !isParenExprwithBinaryExpr(ue.X) {
			l.Err = fmt.Errorf("`not`'s subExpr must be ParenExpr(with BinaryExpr), err at %v", ue.OpPos)
			return false
		}

	case token.SUB:
		basicLit, ok := ue.X.(*ast.BasicLit)
		if !ok || basicLit.Kind != token.INT {
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
		l.Tokens = append(l.Tokens, Token{Typ: lt.Kind, Val: lt.Value})

	case token.STRING:
		l.Tokens = append(l.Tokens, Token{Typ: lt.Kind, Val: strings.Trim(lt.Value, `"`)})

	default:
		return false
	}

	return true
}

// 处理标识符或布尔值
func (l *Lexer) handleIdent(it *ast.Ident) (isValid bool) {
	name := it.Name
	if name == "true" || name == "false" {
		boolVal, _ := strconv.ParseBool(name)
		// bool 值使用 ILLEGAL 特殊标识
		l.Tokens = append(l.Tokens, Token{Typ: token.ILLEGAL, IsBoolean: true, BoolVal: boolVal})
	} else {
		l.Tokens = append(l.Tokens, Token{Typ: token.IDENT, Val: name})
	}

	return true
}

// 生成无效 token 的错误信息
func invalidTokenError(t token.Token, pos token.Pos) error {
	return fmt.Errorf("invalid token(%q) at position(%v)", t, pos)
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
