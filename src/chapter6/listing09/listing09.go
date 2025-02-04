// This sample program demonstrates how to create race
// conditions in our programs. We don't want to do this.
package main

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	// counter is a variable incremented by all goroutines.
	counter int

	// wg is used to wait for the program to finish.
	wg sync.WaitGroup
)

// main is the entry point for all Go programs.
func main() {
	// Add a count of two, one for each goroutine.
	wg.Add(2)

	// Create two goroutines.
	go incCounter(1)
	go incCounter(2)

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

// incCounter increments the package level counter variable.
func incCounter(id int) {
	// Schedule the call to Done to tell main we are done.
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// Capture the value of Counter.
		value := counter

		// Yield the thread and be placed back in queue.
		// 用于将goroutine从当前线程退出，给其他狗routine运行的机会
		// 目的是强制调度器切换两个goroutine，以便让竞争状态的效果变得更明显
		// 这一句很重要。没有这一句，结果是4
		runtime.Gosched()

		// Increment our local value of Counter.
		value++

		// Store the value back into Counter.
		counter = value
	}
}
