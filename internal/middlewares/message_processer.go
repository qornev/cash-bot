package middlewares

import "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"

type MessageProcesser interface {
	IncomingMessage(msg messages.Message, info *messages.CommandInfo) error
}

type MessageProcesserFunc func(msg messages.Message, info *messages.CommandInfo) error

func (f MessageProcesserFunc) IncomingMessage(msg messages.Message, info *messages.CommandInfo) error {
	return f(msg, info)
}

func NewMessageProcesser(model MessageProcesser) MessageProcesser {
	return MessageProcesserFunc(func(msg messages.Message, info *messages.CommandInfo) error {
		return model.IncomingMessage(msg, info)
	})
}
