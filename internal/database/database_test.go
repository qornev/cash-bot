package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

func Test_New_ShouldCreateDatabaseObjects(t *testing.T) {
	_, err := New()
	assert.NoError(t, err)
}

func Test_Add_ShouldAddObjectsToDatabase(t *testing.T) {
	db, _ := New()

	cons := &messages.Consumption{
		Amount:   123.45,
		Category: "еда",
		Date:     time.Now().Unix(),
	}

	err := db.Add(1234, cons)
	assert.NoError(t, err)

	err = db.Add(1234, cons)
	assert.NoError(t, err)

	assert.Equal(t, db.data[1234][0], cons)
	assert.Equal(t, db.data[1234][1], cons)
}

func Test_Get_ShouldReturnSameObjects(t *testing.T) {
	db, _ := New()

	cons := &messages.Consumption{
		Amount:   123.45,
		Category: "еда",
		Date:     time.Now().Unix(),
	}

	db.Add(1234, cons)
	db.Add(1234, cons)

	res, err := db.Get(1234)
	assert.NoError(t, err)

	assert.Equal(t, res[0], cons)
	assert.Equal(t, res[1], cons)
}
