package tg

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/middlewares"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"go.uber.org/zap"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(ctx context.Context, text string, userID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "send message to user")
	defer span.Finish()

	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

var currencyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(converter.RUB, converter.RUB),
		tgbotapi.NewInlineKeyboardButtonData(converter.USD, converter.USD),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(converter.EUR, converter.EUR),
		tgbotapi.NewInlineKeyboardButtonData(converter.CNY, converter.CNY),
	),
)

var MarkupNotExistError = errors.New("keyboard markup not exist")

func (c *Client) SendMessageWithKeyboard(ctx context.Context, text string, keyboardMarkup string, userID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "send keyboard message to user")
	defer span.Finish()

	msg := tgbotapi.NewMessage(userID, text)
	switch keyboardMarkup {
	case "currency":
		msg.ReplyMarkup = currencyKeyboard
	default:
		return MarkupNotExistError
	}

	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Model, clbModel *callbacks.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5 // set timeout for 5s for fast graceful end. prev value 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("listening for messages...")

	for update := range updates {
		if update.Message != nil { // If we got a message
			messageProcesser := middlewares.NewMessageProcesser(msgModel)
			messageProcesser = middlewares.LoggingMiddleware(messageProcesser)
			messageProcesser = middlewares.MetricMiddleware(messageProcesser)
			messageProcesser = middlewares.TracingMiddleware(messageProcesser)

			err := messageProcesser.IncomingMessage(
				messages.Message{
					Text:   update.Message.Text,
					UserID: update.Message.From.ID,
				},
				&messages.CommandInfo{
					Command: messages.Unknown,
				},
			)
			if err != nil {
				logger.Error(
					"error processing message",
					zap.Int64("user_id", update.Message.From.ID),
					zap.String("user_input", update.Message.Text),
				)
			}
		} else if update.CallbackQuery != nil { // If we got a callback
			logger.Info(
				"processing callback",
				zap.Int64("user_id", update.CallbackQuery.From.ID),
				zap.String("user_input", update.CallbackQuery.Data),
			)

			err := clbModel.IncomingCallback(callbacks.Callback{
				Data:   update.CallbackQuery.Data,
				UserID: update.CallbackQuery.From.ID,
				// InlineID: update.CallbackQuery.Message.MessageID,
			})
			if err != nil {
				logger.Error(
					"error processing callback",
					zap.Int64("user_id", update.Message.From.ID),
					zap.String("user_input", update.Message.Text),
				)
			}
		}
	}

	logger.Info("end processing messages...")
}

func (c *Client) AutoListenUpdates(ctx context.Context, wg *sync.WaitGroup, msgModel *messages.Model, clbModel *callbacks.Model) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.ListenUpdates(msgModel, clbModel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		c.client.StopReceivingUpdates()
	}()
}
