package converter_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/converter"
)

// Passed on local but failed in ci. Problem with config file path

// func Test_UpdateRate_ShouldGetNewRates(t *testing.T) {
// 	cfg, _ := config.NewFromFile("../../data/config.yaml")
// 	client := rate.New(cfg)
// 	model := converter.New(client)

// 	err := model.UpdateRate()
// 	assert.NoError(t, err)

// 	assert.NotNil(t, model.GetRate())
// }

// func Test_Exchange_ShouldAnswerWithCorrectValue(t *testing.T) {
// 	cfg, _ := config.NewFromFile("../../data/config.yaml")
// 	client := rate.New(cfg)
// 	model := converter.New(client)

// 	amount, err := model.Exchange(1.0, converter.USD, converter.RUB)
// 	assert.NoError(t, err)

// 	assert.Greater(t, amount, 1.0)
// }

func Test_UpdateRecentRates_ShouldRemoveExtraCurrenciesUserReportsFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	rateUpdater := mocks.NewMockRateUpdater(ctrl)
	rateDB := mocks.NewMockRateManipulator(ctrl)
	userDB := mocks.NewMockUserManipulator(ctrl)
	reportCacher := mocks.NewMockReportCacher(ctrl)

	now := time.Now().Unix()
	rates := &converter.Rates{
		USD: 61.1,
		EUR: 62.1,
		CNY: 10.1,
	}
	users := []domain.User{
		{
			ID:      123,
			Code:    "USD",
			Budget:  nil,
			Updated: now,
		},
		{
			ID:      321,
			Code:    "USD",
			Budget:  nil,
			Updated: now,
		},
		{
			ID:      222,
			Code:    "USD",
			Budget:  nil,
			Updated: now,
		},
	}

	rateUpdater.EXPECT().GetUpdate(gomock.Any(), nil).Return(rates, nil)
	rateDB.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	userDB.EXPECT().GetAllUsers(gomock.Any()).Return(users, nil)

	reportCacher.EXPECT().RemoveFromAll(gomock.Any(), users[0].ID)
	reportCacher.EXPECT().RemoveFromAll(gomock.Any(), users[1].ID)
	reportCacher.EXPECT().RemoveFromAll(gomock.Any(), users[2].ID)

	model := converter.New(rateUpdater, rateDB, userDB, reportCacher)
	model.UpdateRecentRates(context.Background())
}
