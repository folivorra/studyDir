package dirForStudy

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func OneSidedChannel() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()

	go func() {
		for val := range ch {
			fmt.Println(val)
		}
	}()

	time.Sleep(time.Second)
}

func PingPongChannel() {
	pingCh := make(chan int)
	pongCh := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			<-pongCh
			fmt.Println("ping")
			pingCh <- 1

			time.Sleep(200 * time.Millisecond)
		}
	}()

	pongCh <- 1

	go func() {
		for i := 1; i <= 10; i++ {
			<-pingCh
			fmt.Println("pong")
			pongCh <- 1

			time.Sleep(200 * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)
}

func MultichannelInput() {
	rand.Seed(time.Now().UnixNano())
	chans := []chan string{make(chan string), make(chan string), make(chan string)}

	for _, ch := range chans {
		go func(chan string) {
			for i := 1; i <= 10; i++ {
				ch <- fmt.Sprintf("%d", rand.Intn(5)+1)
				time.Sleep(time.Duration(rand.Intn(10)) * 100 * time.Millisecond)
			}
		}(ch)
	}

	counters := make([]int, 3)

	for {
		select {
		case msg := <-chans[0]:
			if counters[0] < 3 {
				fmt.Println("chan 1:", msg)
				counters[0]++
			}
		case msg := <-chans[1]:
			if counters[1] < 3 {
				fmt.Println("chan 2:", msg)
				counters[1]++
			}
		case msg := <-chans[2]:
			if counters[2] < 3 {
				fmt.Println("chan 3:", msg)
				counters[2]++
			}
		case <-time.After(1500 * time.Millisecond):
			fmt.Println("таймаут")
			return
		}

		if counters[0] == 3 && counters[1] == 3 && counters[2] == 3 {
			fmt.Println("все лимиты достигнуты")
			return
		}
	}
}

func TimeoutGoroutine() {
	ch := make(chan int)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Printf("значение пришло\n")
	case <-time.After(2 * time.Second):
		fmt.Println("таймаут")
	}
}

func AnotherFanIn() {
	chs := []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))

	res := anotherMerge(chs)

	for _, ch := range chs {
		go func(chan int) {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Duration(rand.Intn(10)) * 200 * time.Millisecond)
				ch <- rand.Intn(100)
			}
			close(ch)
		}(ch)
	}

	for i := 0; i < 5; i++ {
		fmt.Println(<-res)
	}
}

func anotherMerge(chs []chan int) chan int {
	res := make(chan int)
	wg := &sync.WaitGroup{}

	for _, ch := range chs {
		wg.Add(1)
		go func(chan int) {
			defer wg.Done()
			for {
				val, ok := <-ch
				if !ok {
					return
				}
				res <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}

func AnotherFanOut() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	go func() {
		for i := 1; i < 21; i++ {
			ch <- i
		}
		close(ch)
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go anotherWorker(i, ch, wg)
	}

	wg.Wait()
}

func anotherWorker(id int, in <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range in {
		fmt.Printf("worker %d done job %d\n", id+1, val)
		time.Sleep(1 * time.Second)
	}
}

func BufferChan() {
	ch := make(chan int, 3)

	go func() {
		for val := range ch {
			time.Sleep(time.Second)
			fmt.Printf("прочитано %d\n", val)
		}
	}()

	for i := 0; i < 5; i++ {
		ch <- i
	}

	fmt.Scanln()
}

func DoneChannel() {
	randomCh := make(chan int)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case randomCh <- rand.Intn(100):
			case <-done:
				fmt.Println("done closed")
				return
			}
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	for {
		select {
		case val := <-randomCh:
			fmt.Printf("random - %d\n", val)
			time.Sleep(500 * time.Millisecond)
		case <-time.After(3 * time.Second):
			fmt.Println("timeout")
			return
		}
	}
}

func Semaphor() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sem := make(chan struct{}, 3)
	done := make(chan struct{})

	comleted := 0
	total := len(jobs)

	for _, job := range jobs {
		sem <- struct{}{}
		go semWorker(sem, done, job)
	}

	for comleted < total {
		<-done
		comleted++
	}

	fmt.Println("all workers done")
}

func semWorker(sem chan struct{}, done chan struct{}, job int) {
	time.Sleep(1 * time.Second)
	fmt.Println("worker done ", job)
	<-sem
	done <- struct{}{}
}
