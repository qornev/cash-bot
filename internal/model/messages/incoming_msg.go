package messages

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
	SendMessageWithKeyboard(text string, keyboardMarkup string, userID int64) error
}

type ExpenseManipulator interface {
	Add(ctx context.Context, date int64, userID int64, category string, amount float64) error
	Get(ctx context.Context, userID int64) ([]domain.Expense, error)
}

type UserManipulator interface {
	GetCode(ctx context.Context, userID int64) (string, error)
	SetBudget(ctx context.Context, userID int64, budget float64) error
	GetBudget(ctx context.Context, userID int64) (*float64, string, int64, error)
}

type Converter interface {
	Exchange(value float64, from string, to string) (float64, error)
	UpdateHistoricalRates(date *int64) error
	GetHistoricalCodeRate(code string, date int64) (float64, error)
}

type Model struct {
	tgClient  MessageSender
	userDB    UserManipulator
	expenseDB ExpenseManipulator
	converter Converter
}

func New(tgClient MessageSender, userDB UserManipulator, expenseDB ExpenseManipulator, converter Converter) *Model {
	return &Model{
		tgClient:  tgClient,
		userDB:    userDB,
		expenseDB: expenseDB,
		converter: converter,
	}
}

type Message struct {
	Text   string
	UserID int64
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
/currency - изменить валюту
/set_budget 12.3 - установка лимита на месяц
/show_budget - вывод текущего лимита`

// Messages routing
func (s *Model) IncomingMessage(msg Message) error {
	switch {
	case msg.Text == "/start":
		return s.tgClient.SendMessage(greeting, msg.UserID)
	case msg.Text == "/week" || msg.Text == "/month" || msg.Text == "/year":
		return s.sendReport(msg)
	case msg.Text == "/currency":
		return s.tgClient.SendMessageWithKeyboard("Выберите валюту", "currency", msg.UserID)
	case strings.HasPrefix(msg.Text, "/set_budget"):
		return s.setBudget(msg)
	case msg.Text == "/show_budget":
		return s.getBudget(msg.UserID)
	default:
		// If no match with any command - start parse line
		expense, err := parseExpense(msg.Text)
		if err != nil {
			log.Println(msg.UserID, "expense did not parse")
		}

		if expense != nil {
			err := s.addExpense(expense, msg)
			if err != nil {
				return errors.Wrap(err, "can't add expense")
			}
			return s.tgClient.SendMessage("Расход записан:)", msg.UserID)
		}

		return s.tgClient.SendMessage("Неизвестная команда:(", msg.UserID)
	}
}

// Send prepared report with expenses to user
func (s *Model) sendReport(msg Message) error {
	currentTime := time.Now()
	var startTime time.Time
	switch msg.Text {
	case "/week":
		startTime = currentTime.AddDate(0, 0, -int(currentTime.Weekday())) // Start from Monday
	case "/month":
		startTime = currentTime.AddDate(0, 0, 1-currentTime.Day()) // Start from first day in month
	case "/year":
		startTime = currentTime.AddDate(0, 1-int(currentTime.Month()), 1-currentTime.Day()) // Start with first dat in year
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

// Get calculated user sum of expenses from `startTime`
func (s *Model) getReport(startTime time.Time, msg Message) (string, error) {
	list, err := s.getExpenses(msg.UserID)
	if err != nil {
		return "", errors.Wrap(err, "can't get expenses")
	}

	code, err := s.getCode(msg.UserID)
	if err != nil {
		return "", errors.Wrap(err, "can't get current state")
	}

	sum := make(map[string]float64)
	for _, elem := range list {
		if elem.Date > startTime.Unix() {
			if code != converter.RUB {
				rate, err := s.converter.GetHistoricalCodeRate(code, elem.Date)
				switch {
				case err == sql.ErrNoRows:
					if err = s.converter.UpdateHistoricalRates(&elem.Date); err != nil {
						return "", err
					}

					rate, err = s.converter.GetHistoricalCodeRate(code, elem.Date)
					if err != nil {
						return "", err
					}
				case err != nil:
					return "", err
				}

				elem.Amount = elem.Amount / rate
			}
			sum[elem.Category] += elem.Amount
		}
	}

	var report strings.Builder
	for key, value := range sum {
		if err != nil {
			return "", errors.Wrap(err, "can't convert value")
		}
		fmt.Fprintf(&report, "%s - %.2f %s\n", key, value, code)
	}

	return report.String(), nil
}

// Add expense to database with converting
func (s *Model) addExpense(expense *domain.Expense, msg Message) error {
	code, err := s.getCode(msg.UserID)
	if err != nil {
		return errors.Wrap(err, "can't get code state")
	}

	expense.Amount, err = s.converter.Exchange(expense.Amount, code, converter.RUB)
	if err != nil {
		return errors.Wrap(err, "can't convert value")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.expenseDB.Add(ctx, expense.Date, msg.UserID, expense.Category, expense.Amount)
}

// Get list of all user expenses
func (s *Model) getExpenses(userID int64) ([]domain.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	list, err := s.expenseDB.Get(ctx, userID)
	return list, err
}

// Get user currency code
func (s *Model) getCode(userID int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code, err := s.userDB.GetCode(ctx, userID)
	return code, err
}

// Get user month limit (budget)
func (s *Model) getBudget(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	budget, code, _, err := s.userDB.GetBudget(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "can't get user budget")
	}

	if budget == nil {
		return s.tgClient.SendMessage("У вас не задан бюджет на месяц", userID)
	}

	value, err := s.converter.Exchange(*budget, converter.RUB, code)
	if err != nil {
		return err
	}

	return s.tgClient.SendMessage(fmt.Sprintf("Ваш бюджет на месяц: %.2f %s", value, code), userID)
}

var budgetRe = regexp.MustCompile(`\/set_budget ([0-9.]+[0-9]+$)`)

// Set user budget. Will save it to database
func (s *Model) setBudget(msg Message) (err error) {
	matches := budgetRe.FindStringSubmatch(msg.Text)
	if len(matches) < 2 {
		return ErrorIncorrectLine
	}

	budget, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return errors.Wrap(err, "can't conv budget to float")
	}

	code, err := s.getCode(msg.UserID)
	if err != nil {
		return errors.Wrap(err, "can't get code state")
	}

	budget, err = s.converter.Exchange(budget, code, converter.RUB)
	if err != nil {
		return errors.Wrap(err, "can't convert value")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = s.userDB.SetBudget(ctx, msg.UserID, budget)
	if err != nil {
		return errors.Wrap(err, "can't set budget")
	}

	return s.tgClient.SendMessage("Бюджет установлен", msg.UserID)
}

var expenseRe = regexp.MustCompile(`^([0-9.]+) ([а-яА-Яa-zA-Z]+) ?([0-9]{4}-[0-9]{2}-[0-9]{2})?$`)

var ErrorIncorrectLine = errors.New("Incorrect line")

// Parse line with user expense
func parseExpense(text string) (*domain.Expense, error) {
	matches := expenseRe.FindStringSubmatch(text)
	if len(matches) < 4 {
		return nil, ErrorIncorrectLine
	}

	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, errors.Wrap(err, "can't conv amount of expense")
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

	return &domain.Expense{
		Amount:   amount,
		Category: category,
		Date:     date.Unix(),
	}, nil
}
