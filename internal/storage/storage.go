package storage

import (
	"sync"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

type Storage struct {
	mutex *sync.Mutex
	data  map[int64][]*messages.Expense
	state map[int64]string
}

func New() (*Storage, error) {
	return &Storage{
		mutex: &sync.Mutex{},
		data:  make(map[int64][]*messages.Expense),
		state: make(map[int64]string),
	}, nil
}

func (s *Storage) Add(userID int64, expense *messages.Expense) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[userID] = append(s.data[userID], expense)
	return nil
}

func (s *Storage) Get(userID int64) ([]*messages.Expense, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.data[userID], nil
}

func (s *Storage) SetState(userID int64, currency string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.state[userID] = currency
	return nil
}

func (s *Storage) GetState(userID int64) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if val, ok := s.state[userID]; ok {
		return val, nil
	}
	return converter.RUB, nil
}
