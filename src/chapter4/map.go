package main

import "fmt"

func main() {
	colors := map[string]string{
		"AliceBlue":   "#f0f8ff",
		"Coral":       "#ff7f50",
		"DarkGray":    "#a9a9a9",
		"ForestGreen": "#228b22",
	}

	showAll(colors)

	// 改变映射
	removeColor(colors, "Coral")
	removeColor(colors, "Red") // 删除一个不存在的值，无影响

	fmt.Printf("After change: \n")

	showAll(colors)
}

// 打印所有元素
func showAll(colors map[string]string) {
	fmt.Printf("len: %d\n", len(colors))

	for key, value := range colors {
		fmt.Printf("\tKey: %s Value: %s\n", key, value)
	}
}

func removeColor(colors map[string]string, color string) {
	delete(colors, color)
}
