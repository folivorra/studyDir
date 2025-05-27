package storage

import (
	"encoding/json"
	"fmt"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/model"
	"os"
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

func (s *InMemoryStorage) SaveToFile(path string) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.items, "", "\t")
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (s *InMemoryStorage) LoadFromFile(path string) (err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tempStorage := make(map[int]model.Item)

	if err = json.Unmarshal(data, &tempStorage); err != nil {
		return err
	}

	s.mu.Lock()
	s.items = tempStorage
	s.mu.Unlock()

	return nil
}
