package main

import (
	"fmt"
	"strconv"
	"sync"
)

// Требуется исправить/переписать код до рабочего состояния, который представлен на экране
func main() {
	var wc sync.WaitGroup
	m := make(chan string, 5)
	for i := 0; i < 5; i++ {
		wc.Add(1)
		go func(mm chan<- string, i int, group *sync.WaitGroup) {
			defer group.Done()
			mm <- fmt.Sprintf("Goroutine %s", strconv.Itoa(i))
		}(m, i, &wc)
	}

	wc.Wait()
	close(m)
	for q := range m {
		fmt.Println(q)
	}

}
