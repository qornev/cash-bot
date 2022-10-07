package callbacks_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/model"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
)

func Test_IncomingCallback_ShouldChangeCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	stater := mocks.NewMockStateManipulator(ctrl)
	model := callbacks.New(sender, stater)

	var userID int64 = 1234

	sender.EXPECT().SendMessage("Валюта изменена на USD", userID)
	stater.EXPECT().SetState(userID, "USD")

	err := model.IncomingCallback(callbacks.Callback{
		Data:   "USD",
		UserID: userID,
	})

	assert.NoError(t, err)
}
