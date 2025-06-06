package main

import (
	"fmt"
	"runtime"
	"time"
)

type MyStruct struct {
	A int
	B float64
	C [30]string
}

func main() {
	BlackBox()
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
