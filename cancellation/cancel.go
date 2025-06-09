package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, cancel chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-cancel:
			fmt.Printf("Воркер %d завершает свою работу по сигналу.\n", id)
			return
		default:
			fmt.Printf("Воркер %d выполняет работу.\n", id)
			time.Sleep(400 * time.Millisecond)
		}
	}
}

func main() {
	cancel := make(chan struct{})
	wg := new(sync.WaitGroup)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, cancel, wg)
	}

	time.Sleep(5 * time.Second)
	close(cancel)
	wg.Wait()

	fmt.Println("Все воркеры закончили работу.")
}
