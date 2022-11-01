package middlewares

import (
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"go.uber.org/zap"
)

func LoggingMiddleware(next MessageProcesser) MessageProcesser {
	return MessageProcesserFunc(func(msg messages.Message, info *messages.CommandInfo) error {
		logger.Info(
			"processing message",
			zap.Int64("user_id", msg.UserID),
			zap.String("user_input", msg.Text),
		)

		err := next.IncomingMessage(msg, info)

		logger.Info(
			"message processing complete",
			zap.Int64("user_id", msg.UserID),
			zap.String("command", info.Command),
		)
		return err
	})
}
