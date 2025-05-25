package dirForStudy

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

func ErrorGroup() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	group, ctx := errgroup.WithContext(context.Background())
	// функция создает группу и контекст для этой группы, который отменится если кто то из группы завершится с ошибкой

	for i := 0; i < 5; i++ {
		i := i
		// защита от замыкания
		group.Go(func() error {
			select {
			case <-time.After(time.Second * time.Duration(rand.Intn(5)+1)):
				if i == rand.Intn(5) {
					return fmt.Errorf("worker %d finished with error", i)
				}
				// имитируем случайную ошибку и возвращаем ее -> отменяем контекст группы

				fmt.Printf("worker %d done his job\n", i)
				// если ошибки не было значит воркер выполнил свою работу
				return nil
			case <-ctx.Done():
				// когда контекст группы отменяется - этот канал закрывается и мы возвращаем ошибку завершения группы
				return ctx.Err()
			}
		})
	}

	// ждем выполнение всех горутин в группе и если ошибка ненулевая, то выводим ее, иначе все ок
	if err := group.Wait(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("worker finished without error")
	}
}

// атомики это примитивы синхронизации обеспечивающие атомарные операции с данными (чтение/запись),
// они менее прожорливые чем мьютексы, но используются в простых структурах где необходимо
// менять лишь одно поле (например, счетчик)

type Counter struct {
	value int64
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

func AtomicGoroutine() {
	wg := &sync.WaitGroup{}
	c := &Counter{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("counter %d\n", c.Get())
}

// cond это примитив синхронизации который позволяет усыплять,
// а затем будить горутины по сигналу

func CondExample() {
	cond := sync.NewCond(&sync.Mutex{})

	go func() {
		cond.L.Lock()
		cond.Wait()
		fmt.Println("горутина запустилась!")
		cond.L.Unlock()
	}()

	time.Sleep(time.Second)

	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
} // самый простой вариант использования, также изучил вариант с моделью производитель/потребитель
