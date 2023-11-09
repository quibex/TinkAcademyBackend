package main

import (
	"fmt"
	"sync"
)

func main() {
	wc := &sync.WaitGroup{}
	ch := make(chan string)
	for i := 0; i < 5; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()
			ch <- fmt.Sprintf("Goroutine %d", i)
		}(i)
	}

	go func() {
		wc.Wait()
		close(ch)
	}()

	for q := range ch {
		fmt.Println(q)
	}
}
