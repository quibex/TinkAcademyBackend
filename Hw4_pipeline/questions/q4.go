package main

import "fmt"

// Что выведет данный код?
func main() {
	c := make(chan int, 1)
	c <- 5
	close(c)
	fmt.Println(<-c)
}
