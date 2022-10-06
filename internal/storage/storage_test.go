package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

func Test_New_ShouldCreateStorageObject(t *testing.T) {
	_, err := New()
	assert.NoError(t, err)
}

func Test_Add_ShouldAddConsumptionToStorage(t *testing.T) {
	storage, _ := New()

	cons := &messages.Consumption{
		Amount:   123.45,
		Category: "еда",
		Date:     time.Now().Unix(),
	}

	err := storage.Add(1234, cons)
	assert.NoError(t, err)

	err = storage.Add(1234, cons)
	assert.NoError(t, err)

	assert.Equal(t, storage.data[1234][0], cons)
	assert.Equal(t, storage.data[1234][1], cons)
}

func Test_Get_ShouldReturnSameConsumption(t *testing.T) {
	storage, _ := New()

	cons := &messages.Consumption{
		Amount:   123.45,
		Category: "еда",
		Date:     time.Now().Unix(),
	}

	err := storage.Add(1234, cons)
	assert.NoError(t, err)
	err = storage.Add(1234, cons)
	assert.NoError(t, err)

	res, err := storage.Get(1234)
	assert.NoError(t, err)

	assert.Equal(t, res[0], cons)
	assert.Equal(t, res[1], cons)
}
