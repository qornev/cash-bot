package callbacks

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/base"
)

type Model struct {
	tgClient base.MessageSender
	storage  base.StateManipulator
}

func New(tgClient base.MessageSender, storage base.StateManipulator) *Model {
	return &Model{
		tgClient: tgClient,
		storage:  storage,
	}
}

type Callback struct {
	UserID int64
	//	InlineID int64 for updating messages (later?)
	Data string
}

func (s *Model) IncomingCallback(clb Callback) error {
	err := s.storage.SetState(clb.UserID, clb.Data)
	if err != nil {
		return errors.Wrap(err, "can't set currency state")
	}
	return s.tgClient.SendMessage(fmt.Sprintf("Валюта изменена на %s", clb.Data), clb.UserID)
}
