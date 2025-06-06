package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	MemoryClean()
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
