package dirForStudy

import (
	"errors"
	"sync"
)

func PanicFromMap() {
	m := NewPanicMap()

	for i := 0; i < 100; i++ {
		go func() {
			m.Add("hello", 1)
			m.Get("hello")
		}()
	}
}

func WithoutPanicFromMap() {
	m := NewMapRWMutex()

	for i := 0; i < 100; i++ {
		go func() {
			m.Add("hello", 1)
			m.Get("hello")
		}()
	}
}

type PanicMap struct {
	m map[string]int
}

func NewPanicMap() *PanicMap {
	return &PanicMap{m: make(map[string]int)}
}

func (pm *PanicMap) Add(key string, value int) {
	pm.m[key] = value
}

func (pm *PanicMap) Get(key string) (value int, err error) {
	value, ok := pm.m[key]
	if !ok {
		err = errors.New("key not found")
	}
	return value, err
}

type MapRWMutex struct {
	m  map[string]int
	mu sync.RWMutex
}

// ttl - duration
// serel - redis

// get and set http handlers (clear)

// перезаписывать данные спустя 5 секунд
// сервисы должны завершаться в рамках приоритета

func NewMapRWMutex() *MapRWMutex {
	return &MapRWMutex{m: make(map[string]int)}
}

func (mm *MapRWMutex) Add(key string, value int) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.m[key] = value
}

func (mm *MapRWMutex) Get(key string) (value int, err error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	value, ok := mm.m[key]
	if !ok {
		err = errors.New("key not found")
	}
	return value, err
}
