package storage

import (
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

type Storage struct {
	data  map[int64][]*messages.Consumption
	state map[int64]string
}

func New() (*Storage, error) {
	return &Storage{
		data:  make(map[int64][]*messages.Consumption),
		state: make(map[int64]string),
	}, nil
}

func (s *Storage) Add(userID int64, consumtion *messages.Consumption) error {
	s.data[userID] = append(s.data[userID], consumtion)
	return nil
}

func (s *Storage) Get(userID int64) ([]*messages.Consumption, error) {
	return s.data[userID], nil
}

func (s *Storage) SetState(userID int64, currency string) error {
	s.state[userID] = currency
	return nil
}

func (s *Storage) GetState(userID int64) (string, error) {
	return s.state[userID], nil
}
