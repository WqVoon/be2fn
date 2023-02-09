<h1 align="center"> be2fn </h1>
<p align="center">将布尔表达式编译成纯闭包函数，一次编译，多次执行</p>
<p align="center"><a href="#例子">例子</a> | <a href="#说明">说明</a> | <a href="#原理">原理</a> | <a href="#其他">其他</a></p>

# 例子

```go
package main

import (
	"fmt"

	"github.com/wqvoon/be2fn"
)

func main() {
	// 给 Compile 函数传递一个布尔表达式，编译通过时返回函数 testfn，否则返回错误原因 err
	testfn, err := be2fn.Compile("(val > 0 && val < 10) || (val > -10 && val < -1)")
	if err != nil {
		panic(err)
	}

	// 一个正常的布尔表达式函数，用于验证 Compile 生成的函数的正确性
	trueFn := func(val int) bool {
		return (val > 0 && val < 10) || (val > -10 && val < -1)
	}

	for _, val := range []int{0, 1, 2, 10, -1, -2, -10} {
		// 编译出的函数接受 be2fn.Kv 作为参数（本质上是 map[string]interface{}），
		// key 是变量名，value 是对应的值，变量支持多个，类型支持数字、字符串、布尔值，
		// 变量的类型根据二元表达式中的常量在编译期确定
		testRet, err := testfn(be2fn.Kv{"val": val})
		if err != nil {
			panic(err)
		}

		trueRet := trueFn(val)

		if testRet != trueRet {
			panic(fmt.Sprintf("faild with val(%d)", val))
		}
	}
	// 如果执行到这里，说明编译的函数返回的结果正确
	fmt.Println("done")
}

```

# 说明

- 支持的操作符有 `&&`，`||`，`!`，`==`，`!=`，`<`，`<=`，`>`，`>=`
- 支持内建函数 `in` 用于判断变量是否在整数切片或字符串切片中，如 `in(a, []int{1,2,3})`
- `!` 操作符对应的子表达式必须是包含二元表达式的括号表达式或函数调用，如 `!(a == 1)` 或 `!in(a, []int{1,2,3})`
- 最小的可比较单元是二元表达式，所以不支持 `a`，`a && b` ，需要用 `a == true`， `a==true && b==true`
- 二元表达式的操作数必须一个是常量一个是变量，变量的类型根据常量在编译期确定
- 常量和变量支持数字、字符串、布尔三种类型

# 原理

```go
package main

import "fmt"

type Unit func(int) bool

func getResult(val int) bool {
	return (val > 0 && val < 10) || (val > -10 && val < -1)
}

func main() {
	// `(val > 0 && val < 10) || (val > -10 && val < -1)` 表达式根据优先级解析成 AST，
	// 然后进行后序遍历，可以得到如下的逆波兰表达式：
	// `val 0 > val 10 < && val 10 - > val 1 - < && ||`

	// 使用逆波兰表达式的计算方式对上面的表达式进行计算，对每个二元表达式操作符定义一个操作函数
	// 计算过程中每次遇到操作符时就调用对应的函数生成闭包
	u1 := LargeThan(0)
	u2 := LessThan(10)
	// 与或非逻辑的操作函数和普通的二元表达式不同，它们不接受操作数，接受普通二元表达式生成的闭包并产生新的闭包
	u3 := And(u1, u2)

	u4 := LargeThan(-10)
	u5 := LessThan(-1)
	u6 := And(u4, u5)
	// 最后一个闭包函数就是整体表达式对应的函数，调用它时等价于执行了表达式
	u7 := Or(u3, u6)

	for _, val := range []int{0, 1, 2, 10, -1, -2, -10} {
		if u7(val) != getResult(val) {
			panic(fmt.Sprintf("faild with val(%d)", val))
		}
	}
	fmt.Println("done")
}

func LargeThan(y int) Unit {
	return func(val int) bool {
		return val > y
	}
}

func LessThan(y int) Unit {
	return func(val int) bool {
		return val < y
	}
}

func And(x, y Unit) Unit {
	return func(val int) bool {
		return x(val) && y(val)
	}
}

func Or(x, y Unit) Unit {
	return func(val int) bool {
		return x(val) || y(val)
	}
}

```

# 其他

- 技术不佳，欢迎提交 issue 和 PR 进行交流 :-)