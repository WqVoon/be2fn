package be2fn

import (
	"errors"
	"go/token"
)

// 一个可以被执行并获取结果的函数
type Unit func(Kv) (bool, error)

// shortcut，如果执行 Unit 时遇到错误，那么返回 false，否则直接返回 Unit 的返回值
func (u Unit) GetBool(vars Kv) bool {
	ret, err := u(vars)
	if err != nil {
		return false
	}
	return ret
}

// &&
func And(x, y Unit) Unit {
	return func(vals Kv) (bool, error) {
		xVal, xErr := x(vals)
		if xErr != nil {
			return false, xErr
		}

		yVal, yErr := y(vals)
		if yErr != nil {
			return false, yErr
		}
		return (xVal && yVal), nil
	}
}

// ||
func Or(x, y Unit) Unit {
	return func(vals Kv) (bool, error) {
		xVal, xErr := x(vals)
		if xErr != nil {
			return false, xErr
		}

		yVal, yErr := y(vals)
		if yErr != nil {
			return false, yErr
		}
		return (xVal || yVal), nil
	}
}

// !
func Not(x Unit) Unit {
	return func(vals Kv) (bool, error) {
		xVal, xErr := x(vals)
		if xErr != nil {
			return false, xErr
		}
		return !xVal, nil
	}
}

// 二元运算符函数集，每个运算符都包含这些函数
type OperatorFuncs struct {
	// 变量比较整数
	VarToInt func(varname string, val int) Unit
	// 整数比较变量
	IntToVar func(val int, varname string) Unit

	// 变量比较字符串
	VarToStr func(varname string, val string) Unit
	// 字符串比较变量
	StrToVar func(val string, varname string) Unit

	// 变量比较布尔值
	VarToBool func(varname string, val bool) Unit
	// 布尔值比较变量
	BoolToVar func(val bool, varname string) Unit
}

// 运算符函数集的默认实现
var DefaultOperatorSet = map[token.Token]OperatorFuncs{
	// ==
	token.EQL: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal == val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val == intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal == val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val == strVal), nil
			}
		},

		VarToBool: func(varname string, val bool) Unit {
			return func(vars Kv) (bool, error) {
				boolVal, err := vars.GetBool(varname)
				if err != nil {
					return false, err
				}
				return (boolVal == val), nil
			}
		},

		BoolToVar: func(val bool, varname string) Unit {
			return func(vars Kv) (bool, error) {
				boolVal, err := vars.GetBool(varname)
				if err != nil {
					return false, err
				}
				return (val == boolVal), nil
			}
		},
	},

	// !=
	token.NEQ: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal != val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val != intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal != val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val != strVal), nil
			}
		},

		VarToBool: func(varname string, val bool) Unit {
			return func(vars Kv) (bool, error) {
				boolVal, err := vars.GetBool(varname)
				if err != nil {
					return false, err
				}
				return (boolVal != val), nil
			}
		},

		BoolToVar: func(val bool, varname string) Unit {
			return func(vars Kv) (bool, error) {
				boolVal, err := vars.GetBool(varname)
				if err != nil {
					return false, err
				}
				return (val != boolVal), nil
			}
		},
	},

	// <
	token.LSS: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal < val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val < intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal < val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val < strVal), nil
			}
		},

		VarToBool: CompareBooleanLeft,

		BoolToVar: CompareBooleanRight,
	},

	// <=
	token.LEQ: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal <= val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val <= intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal <= val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val <= strVal), nil
			}
		},

		VarToBool: CompareBooleanLeft,

		BoolToVar: CompareBooleanRight,
	},

	// >
	token.GTR: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal > val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val > intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal > val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val > strVal), nil
			}
		},

		VarToBool: CompareBooleanLeft,

		BoolToVar: CompareBooleanRight,
	},

	// >=
	token.GEQ: {
		VarToInt: func(varname string, val int) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (intVal >= val), nil
			}
		},

		IntToVar: func(val int, varname string) Unit {
			return func(vars Kv) (bool, error) {
				intVal, err := vars.GetInt(varname)
				if err != nil {
					return false, err
				}
				return (val >= intVal), nil
			}
		},

		VarToStr: func(varname string, val string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (strVal >= val), nil
			}
		},

		StrToVar: func(val string, varname string) Unit {
			return func(vars Kv) (bool, error) {
				strVal, err := vars.GetString(varname)
				if err != nil {
					return false, err
				}
				return (val >= strVal), nil
			}
		},

		VarToBool: CompareBooleanLeft,

		BoolToVar: CompareBooleanRight,
	},
}

// 布尔值无法比较大小，所以直接返回错误
func CompareBooleanLeft(varname string, val bool) Unit {
	return func(vars Kv) (bool, error) {
		return false, errors.New("boolean values cannot compare numeric sizes")
	}
}

// 布尔值无法比较大小，所以直接返回错误
func CompareBooleanRight(val bool, varname string) Unit {
	return func(vars Kv) (bool, error) {
		return false, errors.New("boolean values cannot compare numeric sizes")
	}
}
