package tg

import (
	"context"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	client.StopReceivingUpdates()
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
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

func (c *Client) SendMessageWithKeyboard(text string, keyboardMarkup string, userID int64) error {
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
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			err := msgModel.IncomingMessage(messages.Message{
				Text:   update.Message.Text,
				UserID: update.Message.From.ID,
			})
			if err != nil {
				log.Println("error processing message:", err)
			}
		} else if update.CallbackQuery != nil { // If we got a callback
			log.Printf("[%s] send callback data %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

			err := clbModel.IncomingCallback(callbacks.Callback{
				Data:   update.CallbackQuery.Data,
				UserID: update.CallbackQuery.From.ID,
				// InlineID: update.CallbackQuery.Message.MessageID,
			})
			if err != nil {
				log.Println("error processing callback:", err)
			}
		}
	}
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
