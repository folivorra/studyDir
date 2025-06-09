package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type MyStruct struct {
	A int
	B float64
	C [30]string
}

func main() {
	//fmt.Println("Без sync.Pool")
	//start := time.Now()
	//runWithoutPool()
	//fmt.Println("Duration:", time.Since(start))
	//
	//sink = nil // очистка
	//
	//fmt.Println("\nС sync.Pool")
	//start = time.Now()
	//runWithPool()
	//fmt.Println("Duration:", time.Since(start))

	//FromGoToOC()

	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)
	//fmt.Printf("Before: StackInuse = %d KB | HeapAlloc = %d KB\n", m.StackInuse/1024, m.HeapAlloc/1024)
	//
	//recurse(0, 8000)
	//
	//runtime.ReadMemStats(&m)
	//fmt.Printf("After:  StackInuse = %d KB | HeapAlloc = %d KB\n\n", m.StackInuse/1024, m.HeapAlloc/1024)
	//
	//runtime.ReadMemStats(&m)
	//fmt.Printf("Before: StackInuse = %d KB | HeapAlloc = %d KB\n", m.StackInuse/1024, m.HeapAlloc/1024)
	//
	//_ = heapAlloc(8000)
	//
	//runtime.ReadMemStats(&m)
	//fmt.Printf("After:  StackInuse = %d KB | HeapAlloc = %d KB\n", m.StackInuse/1024, m.HeapAlloc/1024)
}

func MemoryClean() {
	debug.SetMemoryLimit(1024 * 1024 * 1024) // 1 GB

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("До запуска горутины — HeapAlloc: %d KB\n", memStats.HeapAlloc/1024)

	stopCh := make(chan struct{})

	go func() {
		fmt.Println("Горутина: начинаю выделять ~200 MB в срез")
		bigSlice := make([]byte, 200*1024*1024) // ≈200 MB

		for i := range bigSlice {
			bigSlice[i] = 1
		}
		fmt.Println("Горутина: память выделена")

		<-stopCh

		bigSlice = nil
		fmt.Println("Горутина: получил stop — обнулил ссылку на срез и вызываю Goexit()")

		runtime.Goexit()
	}()

	time.Sleep(1 * time.Second)

	runtime.ReadMemStats(&memStats)
	fmt.Printf("После запуска горутины — HeapAlloc: %d KB\n", memStats.HeapAlloc/1024)

	// каждые 5 секунд проверяем, не пора ли вызвать GC по порогу в 100 MB

	//memoryLimit := uint64(100 * 1024 * 1024) // 100 MB в байтах
	//ticker := time.NewTicker(5 * time.Second)
	//go func() {
	//	for range ticker.C {
	//		runtime.ReadMemStats(&memStats)
	//		// Если текущий HeapAlloc превысил 100 MB, запускаем сборщик мусора
	//		if memStats.HeapAlloc > memoryLimit {
	//			fmt.Printf(
	//				"Фоновая проверка: текущий HeapAlloc = %d KB > %d KB, вызываю runtime.GC()\n",
	//				memStats.HeapAlloc/1024,
	//				memoryLimit/1024,
	//			)
	//			runtime.GC()
	//		}
	//	}
	//}()

	time.Sleep(5 * time.Second)
	fmt.Println("Main: отправляю сигнал stop в горутину")
	close(stopCh) // этот «принудительно» завершит горутину (наш код внутри паник/Goexit)

	time.Sleep(5 * time.Second)

	runtime.ReadMemStats(&memStats)
	fmt.Printf("После завершения горутины — HeapAlloc: %d KB\n", memStats.HeapAlloc/1024)

	fmt.Println("Main: вызываю runtime.GC() после завершения горутины")
	runtime.GC()

	time.Sleep(500 * time.Millisecond)

	runtime.ReadMemStats(&memStats)
	fmt.Printf("После завершения горутины и GC — HeapAlloc: %d KB\n", memStats.HeapAlloc/1024)

	//fmt.Println("Main: сплю ещё 10 секунд, чтобы поработал фоновый GC-тестер")
	//time.Sleep(10 * time.Second)

	fmt.Println("Main: выхожу")
	//ticker.Stop()
}

func BlackBox() {
	done := make(chan struct{})
	var memStats runtime.MemStats

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				runtime.ReadMemStats(&memStats)
				fmt.Println("HeapAlloc =", memStats.HeapAlloc/1024/1024, "MB", "NumGC =", memStats.NumGC)
			case <-done:
				return
			}
		}
	}()

	var S []*MyStruct

	for i := 0; i < 100000000; i++ {
		s := &MyStruct{
			A: i,
			B: float64(i) * 0.33,
			C: [30]string{"agdfg", "bbcv", "casd"},
		}

		if s.A%2 == 0 {
			_ = s.C[1]
		}

		S = append(S, s)

		s = nil

		if len(S) > 1000000 {
			S = nil
		}
	}

	time.Sleep(10 * time.Millisecond)

	fmt.Println("Main: всё закончилось, выхожу")
	close(done)
}

const (
	iterations = 500_000
	bufSize    = 1024
)

var sink []byte

func printStats(phase string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("[%s] Alloc = %d KB | TotalAlloc = %d KB | NumGC = %d | Mallocs = %d\n",
		phase, m.Alloc/1024, m.TotalAlloc/1024, m.NumGC, m.Mallocs)
}

func runWithoutPool() {
	printStats("NoPool-Before")

	for i := 0; i < iterations; i++ {
		buf := make([]byte, bufSize)
		buf[0] = byte(i) // используем буфер, чтобы не удалился
		sink = append(sink, buf...)
	}

	runtime.GC()
	printStats("NoPool-After")
}

func runWithPool() {
	pool := sync.Pool{
		New: func() any {
			b := make([]byte, bufSize)
			return b
		},
	}

	printStats("WithPool-Before")

	for i := 0; i < iterations; i++ {
		buf := pool.Get().([]byte)
		buf[0] = byte(i)
		pool.Put(buf)
	}

	runtime.GC()
	printStats("WithPool-After")

	// sync.Pool реально где то используют активно??
}

func FromGoToOC() {
	var m runtime.MemStats

	for i := 0; i < 100; i++ {
		buf := make([]byte, 50*1024*1024) // выделяем большой блок
		_ = buf[0]                        // используем, чтобы не оптимизировали
		buf = nil                         // отпускаем
		runtime.GC()                      // просим GC собрать мусор

		time.Sleep(time.Millisecond * 300)
		runtime.ReadMemStats(&m)
		fmt.Printf("Iteration %d - HeapInuse = %d MB | HeapReleased = %d MB | Sys = %d MB\n",
			i+1,
			m.HeapInuse/1024/1024,
			m.HeapReleased/1024/1024,
			m.Sys/1024/1024,
		)
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("END - HeapInuse = %d MB | HeapReleased = %d MB | Sys = %d MB\n",
		m.HeapInuse/1024/1024,
		m.HeapReleased/1024/1024,
		m.Sys/1024/1024,
	)
}

type Node struct {
	next *Node
	data [128]byte
}

func heapAlloc(n int) *Node {
	var head *Node
	for i := 0; i < n; i++ {
		head = &Node{
			next: head,
		}
	}
	return head
}

func recurse(depth, max int) {
	var x [128]byte
	_ = x[0]
	if depth >= max {
		return
	}
	recurse(depth+1, max)
}
