package storage

import (
	"context"
	"sync"

	"github.com/pkg/errors"
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

var TimeoutStorageError = errors.New("storage connection timeout")

func (s *Storage) Add(ctx context.Context, userID int64, expense *messages.Expense) error {
	complete := make(chan struct{}, 1)

	go func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.data[userID] = append(s.data[userID], expense)

		complete <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return TimeoutStorageError
	case <-complete:
		return nil
	}
}

func (s *Storage) Get(ctx context.Context, userID int64) ([]*messages.Expense, error) {
	expenses := make(chan []*messages.Expense, 1)

	go func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		expenses <- s.data[userID]
	}()

	select {
	case <-ctx.Done():
		return nil, TimeoutStorageError
	case value := <-expenses:
		return value, nil
	}
}

func (s *Storage) SetState(ctx context.Context, userID int64, currency string) error {
	complete := make(chan struct{}, 1)

	go func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.state[userID] = currency
		complete <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return TimeoutStorageError
	case <-complete:
		return nil
	}
}

func (s *Storage) GetState(ctx context.Context, userID int64) (string, error) {
	currency := make(chan string, 1)

	go func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if val, ok := s.state[userID]; ok {
			currency <- val
		}
		currency <- converter.RUB
	}()

	select {
	case <-ctx.Done():
		return "", TimeoutStorageError
	case value := <-currency:
		return value, nil
	}
}
