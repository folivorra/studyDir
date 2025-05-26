package storage

import (
	"fmt"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/model"
	"sync"
)

type InMemoryStorage struct {
	mu    sync.RWMutex
	items map[int]model.Item
}

var _ Storager = (*InMemoryStorage)(nil)

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		items: make(map[int]model.Item),
	}
}

func (s *InMemoryStorage) AddItem(item model.Item) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[item.ID]; !ok {
		return fmt.Errorf("item already exists")
	}
	s.items[item.ID] = item
	return nil
}

func (s *InMemoryStorage) GetAllItems() (items []model.Item, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return nil, fmt.Errorf("no items")
	}
	items = make([]model.Item, len(s.items))
	for i, item := range s.items {
		items[i] = item
	}
	return items, nil
}

func (s *InMemoryStorage) UpdateItem(item model.Item) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[item.ID]; !ok {
		return fmt.Errorf("item does not exist")
	}
	s.items[item.ID] = item
	return nil
}

func (s *InMemoryStorage) DeleteItem(id int) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return fmt.Errorf("item does not exist")
	}
	delete(s.items, id)
	return nil
}

func (s *InMemoryStorage) GetItem(id int) (item model.Item, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.items[id]; !ok {
		return model.Item{}, fmt.Errorf("item does not exist")
	}
	return s.items[id], nil
}
