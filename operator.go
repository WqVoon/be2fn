package be2fn

// 一个可以被执行并获取结果的函数
type Unit func(Kv) bool

// 运算符函数集
type OperatorSetIface interface {
	// >
	GreaterThenIntLeft(varname string, val int64) Unit
	GreaterThenIntRight(val int64, varname string) Unit

	// >=
	GreaterThenOrEqualToIntLeft(varname string, val int64) Unit
	GreaterThenOrEqualToIntRight(val int64, varname string) Unit

	// <
	LessThenIntLeft(varname string, val int64) Unit
	LessThenIntRight(val int64, varname string) Unit

	// <=
	LessThenOrEqualToIntLeft(varname string, val int64) Unit
	LessThenOrEqualToIntRight(val int64, varname string) Unit

	// ==
	EqualToIntLeft(varname string, val int64) Unit
	EqualToIntRight(val int64, varname string) Unit
	EqualToStringLeft(varname string, val string) Unit
	EqualToStringRight(val string, varname string) Unit

	// !=
	NotEqualToIntLeft(varname string, val int64) Unit
	NotEqualToIntRight(val int64, varname string) Unit
	NotEqualToStringLeft(varname string, val string) Unit
	NotEqualToStringRight(val string, varname string) Unit
}

// &&
func And(x, y Unit) Unit {
	return func(vals Kv) bool {
		return x(vals) && y(vals)
	}
}

// ||
func Or(x, y Unit) Unit {
	return func(vals Kv) bool {
		return x(vals) || y(vals)
	}
}

// !
func Not(x Unit) Unit {
	return func(vals Kv) bool {
		return !x(vals)
	}
}

// 运算符函数集的默认实现
type DefaultOperatorSet struct{}

var _ OperatorSetIface = &DefaultOperatorSet{}

func (*DefaultOperatorSet) GreaterThenIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) GreaterThenIntRight(val int64, varname string) Unit {
	return nil
}

func (*DefaultOperatorSet) GreaterThenOrEqualToIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) GreaterThenOrEqualToIntRight(val int64, varname string) Unit {
	return nil
}

func (*DefaultOperatorSet) LessThenIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) LessThenIntRight(val int64, varname string) Unit {
	return nil
}

func (*DefaultOperatorSet) LessThenOrEqualToIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) LessThenOrEqualToIntRight(val int64, varname string) Unit {
	return nil
}

func (*DefaultOperatorSet) EqualToIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) EqualToIntRight(val int64, varname string) Unit {
	return nil
}
func (*DefaultOperatorSet) EqualToStringLeft(varname string, val string) Unit {
	return nil
}
func (*DefaultOperatorSet) EqualToStringRight(val string, varname string) Unit {
	return nil
}

func (*DefaultOperatorSet) NotEqualToIntLeft(varname string, val int64) Unit {
	return nil
}
func (*DefaultOperatorSet) NotEqualToIntRight(val int64, varname string) Unit {
	return nil
}
func (*DefaultOperatorSet) NotEqualToStringLeft(varname string, val string) Unit {
	return nil
}
func (*DefaultOperatorSet) NotEqualToStringRight(val string, varname string) Unit {
	return nil
}
