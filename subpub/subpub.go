package main

import (
	"fmt"
	"runtime"
	"time"
)

func publish(subs []chan string, msg string) {
	for _, sub := range subs {
		sub <- msg
	}
}

func notifier(id int, in <-chan string) {
	for msg := range in {
		fmt.Printf("Подписчик номер %d получил сообщение: %s\n", id, msg)
	}
}

func main() {
	subs := []chan string{
		make(chan string),
		make(chan string),
		make(chan string),
	}
	// пул сабов

	for id, ch := range subs {
		go notifier(id, ch)
	}
	// запускаем горутины которые будут ждать уведомления и печатать их

	go func() {
		for i := 0; i < len(subs); i++ {
			msg := fmt.Sprintf("Новость %d", i+1)
			publish(subs, msg)
			time.Sleep(time.Second) // для наглядности
		}
		// отправка уведомлений во все каналы

		for _, ch := range subs {
			close(ch)
		}
		// так избегаем утечки горутин, если убрать то горутины будут вечно ждать в range сообщения
	}()

	fmt.Scanln() // лень было делать wg
	fmt.Println(runtime.NumGoroutine())
}
