package callbacks

import (
	"fmt"

	"github.com/pkg/errors"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type StateManipulator interface {
	SetState(userID int64, currency string) error
}

type Model struct {
	tgClient MessageSender
	storage  StateManipulator
}

func New(tgClient MessageSender, storage StateManipulator) *Model {
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
