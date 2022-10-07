package storage

import (
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/base"
)

type Storage struct {
	data  map[int64][]*base.Expense
	state map[int64]string
}

func New() (*Storage, error) {
	return &Storage{
		data:  make(map[int64][]*base.Expense),
		state: make(map[int64]string),
	}, nil
}

func (s *Storage) Add(userID int64, expense *base.Expense) error {
	s.data[userID] = append(s.data[userID], expense)
	return nil
}

func (s *Storage) Get(userID int64) ([]*base.Expense, error) {
	return s.data[userID], nil
}

func (s *Storage) SetState(userID int64, currency string) error {
	s.state[userID] = currency
	return nil
}

func (s *Storage) GetState(userID int64) (string, error) {
	if val, ok := s.state[userID]; ok {
		return val, nil
	}
	return "RUB", nil
}
