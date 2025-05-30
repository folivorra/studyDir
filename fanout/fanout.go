package main

import (
	"fmt"
	"sync"
	"time"
)

func workerIn(job int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Горутина начинает работу, обрабатывает job =", job)
	time.Sleep(100 * time.Millisecond * time.Duration(job)) // имитация работы
	fmt.Println("Горутина завершает работу, обрабатывал job =", job)
}

func main() {
	numJobs := 5
	jobs := make(chan int)
	wg := &sync.WaitGroup{}
	// задаем кол-во работ и канал для их распределения

	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	// заполняем канал обычными интами и закрываем канал на запись

	for job := range jobs {
		wg.Add(1)
		go workerIn(job, wg)
	}
	// для каждой задачи создается отдельная горутина, синхронизируем через wg

	wg.Wait()

	fmt.Println("Горутины закончили работу")
}
