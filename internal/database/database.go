package database

import (
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

type Database struct {
	data map[int64][]*messages.Consumption
}

func New() (*Database, error) {
	return &Database{
		data: make(map[int64][]*messages.Consumption),
	}, nil
}

func (db *Database) Add(userID int64, consumtion *messages.Consumption) error {
	db.data[userID] = append(db.data[userID], consumtion)
	return nil
}

func (db *Database) Remove(userID int64, consumtion *messages.Consumption) error {
	return nil
}
