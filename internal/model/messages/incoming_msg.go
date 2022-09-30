package messages

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type DataManipulator interface {
	Add(userID int64, consumption *Consumption) error
	Remove(userID int64, consumption *Consumption) error
}

type Model struct {
	tgClient MessageSender
	db       DataManipulator
}

func New(tgClient MessageSender, db DataManipulator) *Model {
	return &Model{
		tgClient: tgClient,
		db:       db,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type Consumption struct {
	Amount   float64
	Category string
	Date     time.Time
}

var lineRe = regexp.MustCompile("^([0-9.]+) ([а-яА-Яa-zA-Z]+) ?([0-9]{4}-[0-9]{2}-[0-9]{2})?$")

var errIncorrectLine = errors.New("Incorrect line")

func parseLine(text string) (*Consumption, error) {
	matches := lineRe.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, errIncorrectLine
	}

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, errors.Wrap(err, "can't conv amount of consumption")
	}

	category := matches[2]

	var date time.Time

	if len(matches) < 4 {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", matches[3])
		if err != nil {
			return nil, errors.Wrap(err, "can't parse date")
		}
	}

	return &Consumption{
		Amount:   amount,
		Category: category,
		Date:     date,
	}, nil
}

func (s *Model) IncomingMessage(msg Message) error {
	if msg.Text == "/start" {
		// 		resp := `Бот для учета расходов

		// Добавить трату: <сумма> <категория> <дата*>
		// * - необязательный параметр
		// Пример: 499 интернет 2022-01-01

		// Команды:
		// /start - запуск бота и инструкция
		// /week - недельный отчет
		// /month - месячный отчет
		// /year - годовой отчет`
		resp := "Бот для учета расходов\n\n" +
			"Добавить трату: <сумма> <категория> <дата*>\n" +
			"* - необязательный параметр\n" +
			"Пример: 499 интернет 2022-01-01\n\n" +
			"Команды:\n" +
			"/start - запуск бота и инструкция\n" +
			"/week - недельный отчет\n" +
			"/month - месячный отчет\n" +
			"/year - годовой отчет"

		return s.tgClient.SendMessage(resp, msg.UserID)
	}

	parsed, _ := parseLine(msg.Text)
	if parsed != nil {
		s.db.Add(msg.UserID, parsed)
		return s.tgClient.SendMessage("Расход записан", msg.UserID)
	}

	return s.tgClient.SendMessage("Неизвестная команда", msg.UserID)
}
