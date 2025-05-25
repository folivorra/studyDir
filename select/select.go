package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(3 * time.Second)
		ch2 <- 2
	}()

	timeout := time.After(4 * time.Second)
	// за пределы цикла, чтобы таймер не запускался постоянно заново в for{}

	for {
		select {
		case <-ch1:
			fmt.Println("первый канал отработал")
		case <-ch2:
			fmt.Println("второй канал отработал")
		case <-timeout:
			fmt.Println("таймаут еррор")
			return
		default:
			fmt.Println("ничего нет")
			time.Sleep(1 * time.Second)
		}
	}
}
