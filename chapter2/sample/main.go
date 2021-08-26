package main

import (
	"log"
	"os"

	// _ 表示对包做初始化操作（init函数），但不使用包里的标识符
	_ "chapter2/sample/matchers"

	"chapter2/sample/search"
)

// init is called prior to main.
func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

// main is the entry point for the program.
func main() {
	// Perform the search for the specified term.
	search.Run("president")
}
