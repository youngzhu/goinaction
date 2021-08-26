/*
测试不同的接收者和不同的调用者之间的关系
接收者：值或指针
调用者：值/类型/接口类型的值/接口类型的指针
*/
package main

// 定义一个接口
type Reciver interface {
	Recive(msg string) error
}

// Reciver 的一个默认实现
type defaultReciver struct{}

/*
func (r defaultReciver) Recive(msg string) error {
	return nil
}

// 方法声明使用值类型的接收者
func testValueTypeReciver() {
	// 申明一个指针类型
	pointerType := new(defaultReciver)
	pointerType.Recive("test")

	var valType defaultReciver
	valType.Recive("test")

	// 这种指针赋值方式不对吗？？
	// cannot use &pointerType (type **defaultReciver) as type Reciver in assignment:
	// 	**defaultReciver does not implement Reciver (missing Recive method)
	var intfPointerType1 Reciver = pointerType // 将指针赋值给接口类型
	intfPointerType1.Recive("test")

	// 正常
	var intfPointerType2 Reciver = &valType // 将指针赋值给接口类型
	intfPointerType2.Recive("test")

	var intfValType Reciver = valType // 将值赋值给接口类型
	intfValType.Recive("test")
}
*/

func (r *defaultReciver) Recive(msg string) error {
	return nil
}

// 方法声明使用指针类型的接收者
func testPointerTypeReciver() {
	// 申明一个指针类型
	pointerType := new(defaultReciver)
	pointerType.Recive("test")

	var valType defaultReciver
	valType.Recive("test")

	var intfPointerType1 Reciver = pointerType // 将指针赋值给接口类型
	intfPointerType1.Recive("test")

	var intfPointerType2 Reciver = &valType // 将指针赋值给接口类型
	intfPointerType2.Recive("test")

	// cannot use valType (type defaultReciver) as type Reciver in assignment:
	// defaultReciver does not implement Reciver (Recive method has pointer receiver)
	var intfValType Reciver = valType // 将值赋值给接口类型
	intfValType.Recive("test")
}
