package main

import "fmt"

func main() {
	source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}

	// 长度和容量相等的好处：在append时会分离出一个新数组
	// slice := source[2:3] // 输出Kiwi，改变了原切片
	slice := source[2:3:3] // 输出Banana

	slice = append(slice, "Kiwi")

	fmt.Print(source[3])
}
