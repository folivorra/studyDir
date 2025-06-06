package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Таска %d выполняется\n", job)
		time.Sleep(time.Second) // имитация полезной работы
	}
}

func main() {
	workerCount := 5 // колво рабочих горутин
	jobsCount := 20  // колво задач
	var wg sync.WaitGroup

	jobs := make(chan int) // канал залач, общий для всех горутин

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(jobs, &wg) // запуск воркеров в колве workerCount
	}

	for j := 0; j < jobsCount; j++ {
		jobs <- j // отправка задач в канал
	}
	close(jobs) // канал закрываем чтобы воркеры не ждали после выполнения всех задач еще

	wg.Wait()
}
