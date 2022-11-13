package messages

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/commands"

	//rate_mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/clients/rate"
	//converter_mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/converter"
	mocks "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/mocks/model/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, nil, nil, nil, nil, nil)

	sender.EXPECT().SendMessage(gomock.Any(), "Неизвестная команда:(", int64(123))

	info := CommandInfo{
		Command: commands.Unknown,
	}

	err := model.IncomingMessage(
		Message{
			Text:   "some text",
			UserID: 123,
		},
		&info,
	)

	assert.Equal(t, commands.Unknown, info.Command)
	assert.NoError(t, err)
}

func Test_OnCurrencyCommand_ShouldAnswerWithKeyboardMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, nil, nil, nil, nil, nil)

	sender.EXPECT().SendMessageWithKeyboard(gomock.Any(), "Выберите валюту", "currency", int64(1234))

	info := CommandInfo{
		Command: commands.Unknown,
	}

	err := model.IncomingMessage(
		Message{
			Text:   "/currency",
			UserID: int64(1234),
		},
		&info,
	)

	assert.Equal(t, commands.GetCurrency, info.Command)
	assert.NoError(t, err)
}

func Test_ParseLine_ShouldFillConsumptionFields(t *testing.T) {
	ctx := context.Background()
	line := "123.4 еда 2020-02-02"

	cons, err := parseExpense(ctx, line)

	date, _ := time.Parse("2006-01-02", "2020-02-02")
	assert.NoError(t, err)
	assert.Equal(t, &domain.Expense{
		Amount:   123.4,
		Category: "еда",
		Date:     date.Unix(),
	}, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_NoDataNoPointBadCategory(t *testing.T) {
	ctx := context.Background()
	line := "1234 еДSdа"

	cons, err := parseExpense(ctx, line)

	assert.NoError(t, err)
	assert.Equal(t, float64(1234), cons.Amount)
	assert.Equal(t, "еДSdа", cons.Category)
	assert.Equal(t, time.Now().Round(time.Hour), time.Unix(cons.Date, 0).Round(time.Hour))
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongLine(t *testing.T) {
	ctx := context.Background()
	line := "1.2s34 еДS3dа "

	cons, err := parseExpense(ctx, line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}

func Test_ParseLine_ShouldFillConsumptionFields_WrongDate(t *testing.T) {
	ctx := context.Background()
	line := "123.4 еда 2020-92-92"

	cons, err := parseExpense(ctx, line)

	assert.Error(t, err)
	assert.Nil(t, cons)
}

func Test_onAddExpense_ShouldListExpense(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	userDB := mocks.NewMockUserManipulator(ctrl)
	expenseDB := mocks.NewMockExpenseManipulator(ctrl)
	reportCacher := mocks.NewMockReportCacher(ctrl)
	conver := mocks.NewMockConverter(ctrl)
	model := New(sender, userDB, expenseDB, reportCacher, conver, nil)

	userID := int64(1234)
	dateString := "2022-09-09"
	date, _ := time.Parse("2006-01-02", dateString)
	category := "еда"
	amount := 1234.56

	userDB.EXPECT().GetCode(gomock.Any(), userID).Return(converter.RUB, nil)
	conver.EXPECT().Exchange(gomock.Any(), amount, converter.RUB, converter.RUB).Return(amount, nil)
	reportCacher.EXPECT().RemoveFromAll(gomock.Any(), []int64{userID})
	expenseDB.EXPECT().Add(gomock.Any(), date.Unix(), userID, category, amount)
	sender.EXPECT().SendMessage(gomock.Any(), "Расход записан:)", userID)

	info := CommandInfo{
		Command: commands.Unknown,
	}
	err := model.IncomingMessage(
		Message{
			Text:   fmt.Sprintf("%.2f %s %s", amount, category, dateString),
			UserID: userID,
		},
		&info,
	)

	assert.Equal(t, commands.AddExpense, info.Command)
	assert.NoError(t, err)
}

// func Test_GetReport_ShouldCheckIfUserReportExistInCache(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	userDB := mocks.NewMockUserManipulator(ctrl)
// 	expenseDB := mocks.NewMockExpenseManipulator(ctrl)
// 	reportCache := mocks.NewMockReportCacher(ctrl)

// 	userID := int64(1234)
// 	report := "ddd - 369.00 RUB\n"

// 	reportCache.EXPECT().GetWeekReport(gomock.Any(), userID).Return(report, true)
// 	reportCache.EXPECT().GetMonthReport(gomock.Any(), userID).Return(report, true)
// 	reportCache.EXPECT().GetYearReport(gomock.Any(), userID).Return(report, true)

// 	model := New(nil, userDB, expenseDB, reportCache, nil)
// 	_, err := model.getReportText(context.Background(), Message{
// 		Text:   CommandWeekReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// 	_, err = model.getReportText(context.Background(), Message{
// 		Text:   CommandMonthReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// 	_, err = model.getReportText(context.Background(), Message{
// 		Text:   CommandYearReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// }

// func Test_GetReport_ShouldSaveUserReportToCache(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	userDB := mocks.NewMockUserManipulator(ctrl)
// 	expenseDB := mocks.NewMockExpenseManipulator(ctrl)
// 	reportCache := mocks.NewMockReportCacher(ctrl)

// 	userID := int64(1234)
// 	report := "ddd - 369.00 RUB\n"
// 	now := time.Now().Unix()

// 	expenseDB.EXPECT().Get(gomock.Any(), userID).Return([]domain.Expense{
// 		{
// 			Amount:   123,
// 			Category: "ddd",
// 			Date:     now,
// 		},
// 		{
// 			Amount:   123,
// 			Category: "ddd",
// 			Date:     now - 100,
// 		},
// 		{
// 			Amount:   123,
// 			Category: "ddd",
// 			Date:     now - 1000,
// 		},
// 	}, nil).Times(3)
// 	userDB.EXPECT().GetCode(gomock.Any(), userID).Return(converter.RUB, nil).Times(3)
// 	reportCache.EXPECT().GetWeekReport(gomock.Any(), userID).Return("", false)
// 	reportCache.EXPECT().GetMonthReport(gomock.Any(), userID).Return("", false)
// 	reportCache.EXPECT().GetYearReport(gomock.Any(), userID).Return("", false)

// 	reportCache.EXPECT().SetWeekReport(gomock.Any(), userID, report)
// 	reportCache.EXPECT().SetMonthReport(gomock.Any(), userID, report)
// 	reportCache.EXPECT().SetYearReport(gomock.Any(), userID, report)

// 	model := New(nil, userDB, expenseDB, reportCache, nil)
// 	_, err := model.getReportText(context.Background(), Message{
// 		Text:   CommandWeekReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// 	_, err = model.getReportText(context.Background(), Message{
// 		Text:   CommandMonthReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// 	_, err = model.getReportText(context.Background(), Message{
// 		Text:   CommandYearReport,
// 		UserID: userID,
// 	})
// 	assert.NoError(t, err)
// }

func Test_addExpense_ShouldRemoveUserReportFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	userDB := mocks.NewMockUserManipulator(ctrl)
	expenseDB := mocks.NewMockExpenseManipulator(ctrl)
	reportCache := mocks.NewMockReportCacher(ctrl)
	conver := mocks.NewMockConverter(ctrl)

	userID := int64(1234)
	expense := &domain.Expense{
		Amount:   123,
		Category: "ddd",
		Date:     time.Now().Unix(),
	}

	userDB.EXPECT().GetCode(gomock.Any(), userID).Return(converter.RUB, nil)
	conver.EXPECT().Exchange(gomock.Any(), expense.Amount, converter.RUB, converter.RUB).Return(expense.Amount, nil)
	expenseDB.EXPECT().Add(gomock.Any(), expense.Date, userID, expense.Category, expense.Amount)

	reportCache.EXPECT().RemoveFromAll(gomock.Any(), []int64{userID})

	model := New(nil, userDB, expenseDB, reportCache, conver, nil)
	err := model.addExpense(
		context.Background(),
		expense,
		Message{
			UserID: userID,
			Text:   "ddd",
		},
	)

	assert.NoError(t, err)
}
