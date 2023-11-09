package main

import "fmt"

// Что выведет данный код?
func main() {
	c := make(chan int)
	c <- 5
	close(c)
	fmt.Println(<-c)
}
