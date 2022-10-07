package messages

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/model"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/base"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	sender.EXPECT().SendMessage("Неизвестная команда:(", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnCurrencyCommand_ShouldAnswerWithKeyboardMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, nil)

	sender.EXPECT().SendMessageWithKeyboard("Выберите валюту", "currency", int64(1234))

	err := model.IncomingMessage(Message{
		Text:   "/currency",
		UserID: int64(1234),
	})

	assert.NoError(t, err)
}

func Test_ParseLine_ShouldFillConsumptionFields(t *testing.T) {
	line := "123.4 еда 2020-02-02"

	cons, err := parseLine(line)

	date, _ := time.Parse("2006-01-02", "2020-02-02")
	assert.NoError(t, err)
	assert.Equal(t, &base.Expense{
		Amount:   123.4,
		Category: "еда",
		Date:     date.Unix(),
	}, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_NoDataNoPointBadCategory(t *testing.T) {
	line := "1234 еДSdа"

	cons, err := parseLine(line)

	assert.NoError(t, err)
	assert.Equal(t, float64(1234), cons.Amount)
	assert.Equal(t, "еДSdа", cons.Category)
	assert.Equal(t, time.Now().Round(time.Hour), time.Unix(cons.Date, 0).Round(time.Hour))
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongLine(t *testing.T) {
	line := "1.2s34 еДS3dа "

	cons, err := parseLine(line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongDate(t *testing.T) {
	line := "123.4 еда 2020-92-92"

	cons, err := parseLine(line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}
