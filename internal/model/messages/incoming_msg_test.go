package messages

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"

	//rate_mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/clients/rate"
	//converter_mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/converter"
	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/model/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, nil, nil, nil)

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
	model := New(sender, nil, nil, nil)

	sender.EXPECT().SendMessageWithKeyboard("Выберите валюту", "currency", int64(1234))

	err := model.IncomingMessage(Message{
		Text:   "/currency",
		UserID: int64(1234),
	})

	assert.NoError(t, err)
}

func Test_ParseLine_ShouldFillConsumptionFields(t *testing.T) {
	line := "123.4 еда 2020-02-02"

	cons, err := parseExpense(line)

	date, _ := time.Parse("2006-01-02", "2020-02-02")
	assert.NoError(t, err)
	assert.Equal(t, &domain.Expense{
		Amount:   123.4,
		Category: "еда",
		Date:     date.Unix(),
	}, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_NoDataNoPointBadCategory(t *testing.T) {
	line := "1234 еДSdа"

	cons, err := parseExpense(line)

	assert.NoError(t, err)
	assert.Equal(t, float64(1234), cons.Amount)
	assert.Equal(t, "еДSdа", cons.Category)
	assert.Equal(t, time.Now().Round(time.Hour), time.Unix(cons.Date, 0).Round(time.Hour))
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongLine(t *testing.T) {
	line := "1.2s34 еДS3dа "

	cons, err := parseExpense(line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongDate(t *testing.T) {
	line := "123.4 еда 2020-92-92"

	cons, err := parseExpense(line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}

func Test_onAddExpense_ShouldListExpense(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	userDB := mocks.NewMockUserManipulator(ctrl)
	expenseDB := mocks.NewMockExpenseManipulator(ctrl)
	conver := mocks.NewMockConverter(ctrl)
	model := New(sender, userDB, expenseDB, conver)

	userID := int64(1234)
	dateString := "2022-09-09"
	date, _ := time.Parse("2006-01-02", dateString)
	category := "еда"
	amount := 1234.56

	userDB.EXPECT().GetCode(gomock.Any(), userID).Return(converter.RUB, nil)
	conver.EXPECT().Exchange(amount, converter.RUB, converter.RUB).Return(amount, nil)
	expenseDB.EXPECT().Add(gomock.Any(), date.Unix(), userID, category, amount)
	sender.EXPECT().SendMessage("Расход записан:)", userID)

	err := model.IncomingMessage(Message{
		Text:   fmt.Sprintf("%.2f %s %s", amount, category, dateString),
		UserID: userID,
	})
	assert.NoError(t, err)
}
