package messages

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type DataManipulator interface {
	Add(userID int64, consumption *Consumption) error
	Get(userID int64) ([]*Consumption, error)
}

type StateManipulator interface {
	SetState(userID int64, currency string) error
	GetState(userID int64) (string, error)
}

type StorageManipulator interface {
	DataManipulator
	StateManipulator
}

type Model struct {
	tgClient MessageSender
	storage  StorageManipulator
}

func New(tgClient MessageSender, storage StorageManipulator) *Model {
	return &Model{
		tgClient: tgClient,
		storage:  storage,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type CurrencyState struct {
	Currency string
	UserID   int64
}

type Consumption struct {
	Amount   float64
	Category string
	Date     int64
}

const greeting = `Бот для учета расходов

Добавить трату: <сумма> <категория> <дата*>
* - необязательный параметр
Пример: 499.99 интернет 2022-01-01

Команды:
/start - запуск бота и инструкция
/week - недельный отчет
/month - месячный отчет
/year - годовой отчет
/currency - изменить валюту`

func (s *Model) IncomingMessage(msg Message) error {
	switch msg.Text {
	case "/start":
		return s.tgClient.SendMessage(greeting, msg.UserID)
	case "/week", "/month", "/year":
		return s.sendReport(msg)
	default:
		// If no match with any command - start parse line
		parsed, err := parseLine(msg.Text)
		if err != nil {
			log.Println(msg.UserID, "consumption did not parse")
		}

		if parsed != nil {
			err := s.storage.Add(msg.UserID, parsed)
			if err != nil {
				return errors.Wrap(err, "consumption did not add to db")
			} else {
				return s.tgClient.SendMessage("Расход записан:)", msg.UserID)
			}
		}

		return s.tgClient.SendMessage("Неизвестная команда:(", msg.UserID)
	}
}

func (s *Model) sendReport(msg Message) error {
	currentTime := time.Now()
	var startTime time.Time
	switch msg.Text {
	case "/week":
		startTime = currentTime.AddDate(0, 0, -int(currentTime.Weekday()))
	case "/month":
		startTime = currentTime.AddDate(0, 0, 1-currentTime.Day())
	case "/year":
		startTime = currentTime.AddDate(0, 1-int(currentTime.Month()), 1-currentTime.Day())
	}
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	report, err := s.getReport(startTime, msg)
	if err != nil {
		log.Println(msg.UserID, "can't get report")
		return s.tgClient.SendMessage("Ошибка вывода отчета", msg.UserID)
	}

	if len(report) == 0 {
		return s.tgClient.SendMessage("Для начала добавьте покупки", msg.UserID)
	}
	return s.tgClient.SendMessage("Отчет:\n"+report, msg.UserID)
}

func (s *Model) getReport(startTime time.Time, msg Message) (string, error) {
	list, err := s.storage.Get(msg.UserID)
	if err != nil {
		return "", errors.Wrap(err, "can't get consumption")
	}

	sum := make(map[string]float64)
	for _, elem := range list {
		if elem.Date > startTime.Unix() {
			sum[elem.Category] += elem.Amount
		}
	}

	var resp strings.Builder
	for key, value := range sum {
		fmt.Fprintf(&resp, "%s - %.2f\n", key, value)
	}

	return resp.String(), nil
}

var lineRe = regexp.MustCompile("^([0-9.]+) ([а-яА-Яa-zA-Z]+) ?([0-9]{4}-[0-9]{2}-[0-9]{2})?$")

var errIncorrectLine = errors.New("Incorrect line")

func parseLine(text string) (*Consumption, error) {
	matches := lineRe.FindStringSubmatch(text)
	if len(matches) < 4 {
		return nil, errIncorrectLine
	}

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, errors.Wrap(err, "can't conv amount of consumption")
	}

	category := matches[2]

	var date time.Time
	if len(matches[3]) == 0 { // If no date
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
		Date:     date.Unix(),
	}, nil
}
