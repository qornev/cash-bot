package messages

import (
	"strings"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/kafka/title"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/commands"
	"go.uber.org/zap"
)

// Messages routing
func (s *Model) IncomingMessage(msg Message, info *CommandInfo) error {
	ctx := info.Context()

	switch {
	case msg.Text == commands.CommandStart:
		info.Command = commands.Start
		return s.tgClient.SendMessage(ctx, commands.Greeting, msg.UserID)

	case msg.Text == commands.CommandWeekReport || msg.Text == commands.CommandMonthReport || msg.Text == commands.CommandYearReport:
		info.Command = commands.CommandReportText(msg.Text)
		err := s.producer.ProduceMessage(title.Reports, msg.UserID, msg.Text)
		if err != nil {
			logger.Error("cannot get user report", zap.Int64("user_id", msg.UserID), zap.Error(err))
			return s.tgClient.SendMessage(ctx, "Ошибка формирования отчета", msg.UserID)
		}
		return nil

	case msg.Text == commands.CommandGetCurrency:
		info.Command = commands.GetCurrency
		return s.tgClient.SendMessageWithKeyboard(ctx, "Выберите валюту", "currency", msg.UserID)

	case strings.HasPrefix(msg.Text, commands.CommandSetBudget):
		info.Command = commands.SetBudget
		err := s.setBudget(ctx, msg)
		if err != nil {
			logger.Error("cannot set user budget", zap.Int64("user_id", msg.UserID), zap.Error(err))
			return s.tgClient.SendMessage(ctx, "Не удалось установить бюджет на месяц", msg.UserID)
		}
		return s.tgClient.SendMessage(ctx, "Бюджет на месяц установлен", msg.UserID)

	case msg.Text == commands.CommandShowBudget:
		info.Command = commands.ShowBudget
		text, err := s.getBudgetText(ctx, msg.UserID)
		if err != nil {
			logger.Error("cannot get user budget", zap.Int64("user_id", msg.UserID), zap.Error(err))
			return s.tgClient.SendMessage(ctx, "Не удалось получить ваш бюджет на месяц", msg.UserID)
		}
		return s.tgClient.SendMessage(ctx, text, msg.UserID)

	default:
		// If no match with any command - start parse line
		expense, err := parseExpense(ctx, msg.Text)
		if err != nil {
			// Not `Error` level cause to low importance of this error
			logger.Info("user expense did not parse", zap.String("user_input", msg.Text), zap.Int64("user_id", msg.UserID))
		}

		if expense != nil {
			info.Command = commands.AddExpense
			err := s.addExpense(ctx, expense, msg)
			if err != nil {
				logger.Error("cannot add user expense", zap.Int64("user_id", msg.UserID), zap.Error(err))
				return s.tgClient.SendMessage(ctx, "Не удалось записать трату", msg.UserID)
			}
			return s.tgClient.SendMessage(ctx, "Расход записан:)", msg.UserID)
		}

		info.Command = commands.Unknown
		return s.tgClient.SendMessage(ctx, "Неизвестная команда:(", msg.UserID)
	}
}
