package base

type MessageSender interface {
	SendMessage(text string, userID int64) error
	SendMessageWithKeyboard(text string, keyboardMarkup string, userID int64) error
}

type DataManipulator interface {
	Add(userID int64, expense *Expense) error
	Get(userID int64) ([]*Expense, error)
}

type StateManipulator interface {
	SetState(userID int64, currency string) error
	GetState(userID int64) (string, error)
}

type StorageManipulator interface {
	DataManipulator
	StateManipulator
}

type Expense struct {
	Amount   float64
	Category string
	Date     int64
}

type CurrencyState struct {
	Currency string
	UserID   int64
}
