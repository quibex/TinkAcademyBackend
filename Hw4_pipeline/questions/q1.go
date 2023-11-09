package main

import (
	"fmt"
	"sync"
)

// Что выведет данный код?

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go log("message 1", &wg)
	go log("message 2", &wg)
	wg.Wait()
}

func log(s string, wg *sync.WaitGroup) {
	fmt.Println(s)
	wg.Done()
}
