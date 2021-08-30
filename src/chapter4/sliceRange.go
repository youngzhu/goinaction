package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}

	for index, value := range slice {
		// 迭代返回的变量是迭代过程中根据切片依次赋值的新变量，所以value的地址总是相同的
		// 要想获取每个元素的地址，可以使用切片变量和索引值，如 &slice[index]
		fmt.Printf("Value: %d ValueAddr: %X ElementAddr: %X\n", value, &value, &slice[index])
	}
}
