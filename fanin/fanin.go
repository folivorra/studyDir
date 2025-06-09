package main

import (
	"fmt"
	"sync"
)

func producer(id int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; i < 3; i++ {
			out <- id + i*10
		}
	}()

	// поставляем в канал входные данные (блокируется до момента прочтения)

	return out
}

func fanIn(input ...<-chan int) <-chan int {
	merged := make(chan int)
	wg := sync.WaitGroup{}

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			merged <- n
		}
		// здесь разблокируется, тк читаем по очередно, затем вычитаем wg
	}

	wg.Add(len(input))
	for _, in := range input {
		go output(in)
	}
	// запуск по количеству входных каналов фунцию для объединения каналов

	go func() {
		wg.Wait()
		close(merged)
	}()
	// оборачиваем в горутину чтобы не блокироваться навсегда в моменте записи в merged, иначе wg не обнулиться никогда

	return merged
}

func main() {
	p1, p2, p3 := producer(1), producer(2), producer(3)

	output := fanIn(p1, p2, p3)

	for value := range output {
		fmt.Println(value)
	}
	// range не блокируется благодаря close в fanIn
}
