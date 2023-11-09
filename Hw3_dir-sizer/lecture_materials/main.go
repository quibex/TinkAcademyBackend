package main

import "fmt"

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR Log: %s \n", err)
		}

		fmt.Println("Hello 1")
	}()
	//var ch chan int
	ch := make(chan int)

	go func() {
		defer close(ch)
		ch <- 1
	}()

	fmt.Println(<-ch)
	close(ch) // panic: close of closed channel
	fmt.Println("Hello 2")
}
