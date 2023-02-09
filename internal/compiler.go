package internal

import (
	"errors"
	"fmt"
)

type Compiler struct {
	lex      *Lexer
	units    []Unit   // 子表达式生成的 Unit
	literals []*Param // 操作数
}

func NewCompiler(l *Lexer) *Compiler {
	return &Compiler{
		lex: l,
	}
}

func (c *Compiler) Compile() (Unit, error) {
	for _, t := range c.lex.Params {
		switch t.Typ {
		case IDENT, INT, STRING, BOOLEAN: // 操作数直接入栈供操作符使用
			c.literals = append(c.literals, t)

		case SUB: // 出现减号说明有负数，取栈顶的一个 literal 做处理
			lastIdx := len(c.literals) - 1
			lastVal := c.literals[lastIdx]
			if lastVal.Typ != INT {
				return nil, errors.New("invalid `-` token")
			}
			lastVal.IntVal = -lastVal.IntVal

		case NOT: // not 逻辑，取栈顶的一个 unit 做处理
			lastIdx := len(c.units) - 1
			if len(c.units) == 0 {
				return nil, errors.New("invalid `!` token")
			}
			c.units[lastIdx] = Not(c.units[lastIdx])

		case LAND: // and 逻辑，取栈顶的两个 unit 做处理
			lastIdx := len(c.units) - 1
			if len(c.units) < 2 {
				fmt.Println("units:", c.units)
				return nil, errors.New("invalid `&&` token")
			}
			c.units[lastIdx-1] = And(c.units[lastIdx-1], c.units[lastIdx])
			c.units = c.units[:lastIdx]

		case LOR: // or 逻辑，取栈顶的两个 unit 做处理
			lastIdx := len(c.units) - 1
			if len(c.units) < 2 {
				return nil, errors.New("invalid `||` token")
			}
			c.units[lastIdx-1] = Or(c.units[lastIdx-1], c.units[lastIdx])
			c.units = c.units[:lastIdx]

		case EQL, NEQ, LSS, LEQ, GTR, GEQ:
			u, err := c.handleOperator(t.Typ)
			if err != nil {
				return nil, err
			}
			c.units = append(c.units, u)

		default: // 剩下的 token 被认为是无效的
			return nil, fmt.Errorf("invalid `%s` token", t.Typ)
		}
	}

	if len(c.units) != 1 || len(c.literals) != 0 { // 最终应该只剩一个 unit，没有多余的 literal
		fmt.Println(c.units, c.literals)
		return nil, errors.New("invalid token sequence")
	}
	return c.units[0], nil
}

// 处理二元运算符
func (c *Compiler) handleOperator(t Token) (Unit, error) {
	if len(c.literals) < 2 {
		return nil, fmt.Errorf("invalid `%s` token", t)
	}

	lastIdx := len(c.literals) - 1
	x, y := c.literals[lastIdx-1], c.literals[lastIdx]
	c.literals = c.literals[:lastIdx-1]
	if x.Typ == IDENT { // x 是变量
		opFuncs := DefaultOperatorSet[t]

		switch y.Typ {
		case BOOLEAN: // y 是布尔值
			return opFuncs.VarToBool(x.Val, y.BoolVal), nil
		case INT: // y 是数字
			return opFuncs.VarToInt(x.Val, y.IntVal), nil
		case STRING: // y 是字符串
			return opFuncs.VarToStr(x.Val, y.Val), nil
		default:
			return nil, fmt.Errorf("invalid `%s` token", t)
		}
	}

	if y.Typ == IDENT { // y 是变量
		opFuncs := DefaultOperatorSet[t]

		switch x.Typ {
		case BOOLEAN: // x 是布尔值
			return opFuncs.BoolToVar(x.BoolVal, y.Val), nil
		case INT: // x 是数字
			return opFuncs.IntToVar(x.IntVal, y.Val), nil
		case STRING: // x 是字符串
			return opFuncs.StrToVar(x.Val, y.Val), nil
		default:
			return nil, fmt.Errorf("invalid `%s` token", t)
		}
	}

	// 不该出现的情况
	return nil, fmt.Errorf("invalid `%s` token", t)
}
