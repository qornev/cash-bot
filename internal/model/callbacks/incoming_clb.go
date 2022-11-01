package callbacks

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

type MessageSender interface {
	SendMessage(ctx context.Context, text string, userID int64) error
}

type UserManipulator interface {
	SetCode(ctx context.Context, userID int64, code string) error
}

type Model struct {
	tgClient MessageSender
	userDB   UserManipulator
}

func New(tgClient MessageSender, userDB UserManipulator) *Model {
	return &Model{
		tgClient: tgClient,
		userDB:   userDB,
	}
}

type Callback struct {
	UserID int64
	//	InlineID int64 for updating messages (later?)
	Data string
}

// Callbacks routing
func (s *Model) IncomingCallback(clb Callback) error {
	ctx := context.Background()

	err := s.setCode(clb.UserID, clb.Data)
	if err != nil {
		logger.Error("cannot set code state", zap.Int64("user_id", clb.UserID), zap.Error(err))
		return s.tgClient.SendMessage(ctx, "Не удалось изменить валюту", clb.UserID)
	}

	return s.tgClient.SendMessage(ctx, fmt.Sprintf("Валюта изменена на %s", clb.Data), clb.UserID)
}

// Set currency with `code` to user
func (s *Model) setCode(userID int64, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.userDB.SetCode(ctx, userID, code)
	return err
}
