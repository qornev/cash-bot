package rate

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
)

// Passed on local but failed in ci. Problem with config file path

func Test_getRequestRate_ShouldAnswerWithRates(t *testing.T) {
	cfg, _ := config.NewFromFile("../../../data/config.yaml")
	service := New(cfg)

	resp, err := service.getRequestRate(context.Background())
	assert.NoError(t, err)

	responseRate, err := parseRates(resp)
	assert.NoError(t, err)

	rate := changeEURBaseToRUB(responseRate)
	fmt.Println(rate)

	assert.NotEmpty(t, rate)
}
