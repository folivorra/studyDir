package dirForStudy

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func MainMerge() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// создаем контекст с отменой, также откладываем освобождение ресурсов

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, os.Interrupt)
	// канал с прерыванием

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	ch5 := make(chan int)
	chanSlice := []chan int{ch1, ch2, ch3, ch4, ch5}
	// создаем каналы и собираем в слайс

	res, errs := merge(ctx, ch1, ch2, ch3, ch4, ch5)
	// отправляем все в мерж

	for _, ch := range chanSlice {
		go generator(ch)
	}
	// генератор посылает в канал псевдослучайные значения

	for {
		select {
		// если получили прерывание - освобождаем ресурсы, отменяем контекст и завершаем работу
		case <-sigCh:
			cancel()
			fmt.Println("program ended with signal")
		// получаем значение из res, если он закрыт и значений нет - выходим
		case n, ok := <-res:
			if !ok {
				return
			}
			fmt.Println("value:", n)
		// получаем ошибки из errs, если он закрыт и значений нет - выходим
		case err, ok := <-errs:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		// если отменили контекст - выходим
		case <-ctx.Done():
			return
		}
	}
}

func merge(ctx context.Context, chans ...<-chan int) (<-chan int, <-chan error) {
	result := make(chan int)
	err := make(chan error)
	wg := &sync.WaitGroup{}
	// создаем output каналы и wg

	for _, ch := range chans {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					err <- fmt.Errorf("panic: %v", r)
				}
			}() // рековер механизм

			for val := range ch {
				if val > 100 {
					err <- fmt.Errorf("error: %v > 100", val)
				}
				// если значение больше ста - пишем в канал ошибку
				select {
				case result <- val:
				// если контекст отменили - выходим
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		// пайплайн / воркер пул для отловки блокировки
		wg.Wait()
		close(result)
		close(err)
	}()
	// ждем когда отработают все горутины, затем закрываем каналы
	// обернули в горутину, чтобы функция merge сразу вернула output без блокировки

	return result, err
}

func generator(ch chan int) <-chan int {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		ch <- rand.Intn(130)
		time.Sleep(time.Second)
	}
	close(ch)
	// генерируем значения и пишем в канал с задержкой в одну секунду
	// затем закрываем канал

	return ch
}
