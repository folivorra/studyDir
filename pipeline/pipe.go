package main

import (
	"fmt"
	"strconv"
)

const countNums = 10

func generator() <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; i < countNums; i++ {
			out <- i
		}
	}()

	return out
}

func square(input <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range input {
			out <- n * n
		}
	}()

	return out
}

func itoa(input <-chan int) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for n := range input {
			out <- strconv.Itoa(n)
		}
	}()

	return out
}

func main() {
	for val := range itoa(square(generator())) {
		fmt.Println(val)
	}
}
