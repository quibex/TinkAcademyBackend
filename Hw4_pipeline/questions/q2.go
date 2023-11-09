package main

// Что выведет данный код?
func main() {
	c := make(chan int)
	close(c)
	c <- 5
}
