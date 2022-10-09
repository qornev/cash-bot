package callbacks

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/model/callbacks"
)

func Test_IncomingCallback_ShouldChangeCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	stater := mocks.NewMockStateManipulator(ctrl)
	model := New(sender, stater)

	var userID int64 = 1234

	sender.EXPECT().SendMessage("Валюта изменена на USD", userID)
	stater.EXPECT().SetState(userID, converter.USD)

	err := model.IncomingCallback(Callback{
		Data:   converter.USD,
		UserID: userID,
	})

	assert.NoError(t, err)
}
