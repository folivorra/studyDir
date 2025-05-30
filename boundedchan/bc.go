package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Сообщение %d получено\n", job)
		time.Sleep(2 * time.Second) // имитация медленной работы
	}
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go worker(jobs, &wg)

	for i := 0; i < 15; i++ {
		fmt.Println("Пытаемся отправить сообщение", i)
		jobs <- i
		fmt.Println("Отправили", i)
	}
	// первые 10 сообщений быстро отправятся в канал, но пока нет блокировки не читаются
	// затем воркер все же начинает медленно читать и освобождать буфер
	close(jobs)

	wg.Wait()
}
