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
	stater := mocks.NewMockUserManipulator(ctrl)
	reportCache := mocks.NewMockReportCacher(ctrl)
	model := New(sender, stater, reportCache)

	var userID int64 = 1234

	reportCache.EXPECT().RemoveFromAll(gomock.Any(), gomock.Any()).Return(nil)

	sender.EXPECT().SendMessage(gomock.Any(), "Валюта изменена на USD", userID)
	stater.EXPECT().SetCode(gomock.Any(), userID, converter.USD)

	err := model.IncomingCallback(Callback{
		Data:   converter.USD,
		UserID: userID,
	})

	assert.NoError(t, err)
}

func Test_IncomingCallback_ShouldRemoveUserReportFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	stater := mocks.NewMockUserManipulator(ctrl)
	reportCache := mocks.NewMockReportCacher(ctrl)

	var userID int64 = 1234

	sender.EXPECT().SendMessage(gomock.Any(), gomock.Any(), userID)
	stater.EXPECT().SetCode(gomock.Any(), userID, gomock.Any())

	reportCache.EXPECT().RemoveFromAll(gomock.Any(), []int64{userID})

	model := New(sender, stater, reportCache)
	err := model.IncomingCallback(Callback{
		Data:   converter.USD,
		UserID: userID,
	})

	assert.NoError(t, err)
}
