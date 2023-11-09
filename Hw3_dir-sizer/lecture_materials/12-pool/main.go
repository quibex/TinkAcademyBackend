package main

import (
	"fmt"
	"sync"
)

const WorkerPool = 2

func worker(s string) {
	// hard work
	fmt.Printf("Hello %s \n", s)
	panic("kek")
}

func main() {
	a := []string{"Максим", "Люся", "Миша", "Илья", "Петя"}
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, WorkerPool)

	for _, v := range a {
		ch <- struct{}{}
		wg.Add(1)

		go func(v string) {
			defer func() {
				<-ch
				wg.Done()
			}()

			worker(v)
		}(v)
	}

	wg.Wait()
}
