package messages

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/commands"
)

type MessageSender interface {
	SendMessage(ctx context.Context, text string, userID int64) error
	SendMessageWithKeyboard(ctx context.Context, text string, keyboardMarkup string, userID int64) error
}

type ExpenseManipulator interface {
	Add(ctx context.Context, date int64, userID int64, category string, amount float64) error
}

type UserManipulator interface {
	GetCode(ctx context.Context, userID int64) (string, error)
	SetBudget(ctx context.Context, userID int64, budget float64) error
	GetBudget(ctx context.Context, userID int64) (*float64, string, int64, error)
}

type ReportCacher interface {
	RemoveFromAll(ctx context.Context, key []int64) error
}

type Converter interface {
	Exchange(ctx context.Context, value float64, from string, to string) (float64, error)
}

type Producer interface {
	ProduceMessage(topic string, userID int64, text string) error
}

type Model struct {
	tgClient    MessageSender
	userDB      UserManipulator
	expenseDB   ExpenseManipulator
	reportCache ReportCacher
	converter   Converter
	producer    Producer
}

func New(tgClient MessageSender, userDB UserManipulator, expenseDB ExpenseManipulator, reportCache ReportCacher, converter Converter, producer Producer) *Model {
	return &Model{
		tgClient:    tgClient,
		userDB:      userDB,
		expenseDB:   expenseDB,
		reportCache: reportCache,
		converter:   converter,
		producer:    producer,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type CommandInfo struct {
	Command string
	ctx     context.Context
}

func (c *CommandInfo) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *CommandInfo) WithContext(ctx context.Context) *CommandInfo {
	if ctx == nil {
		panic("nil context")
	}
	c2 := new(CommandInfo)
	*c2 = *c
	c2.ctx = ctx
	return c2
}

// Add expense to database with converting
func (s *Model) addExpense(ctx context.Context, expense *domain.Expense, msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add expense")
	span.SetTag("command", commands.AddExpense)
	defer span.Finish()

	if err := s.reportCache.RemoveFromAll(ctx, []int64{msg.UserID}); err != nil {
		return errors.Wrap(err, "cannot remove user report from cache")
	}

	code, err := s.getCode(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "can't get code state")
	}

	expense.Amount, err = s.converter.Exchange(ctx, expense.Amount, code, converter.RUB)
	if err != nil {
		return errors.Wrap(err, "can't convert value")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.expenseDB.Add(ctx, expense.Date, msg.UserID, expense.Category, expense.Amount)
}

// Get user currency code
func (s *Model) getCode(ctx context.Context, userID int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get code")
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	code, err := s.userDB.GetCode(ctx, userID)
	return code, err
}

// Get user month limit (budget)
func (s *Model) getBudgetText(ctx context.Context, userID int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get budget text")
	span.SetTag("command", commands.ShowBudget)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	budget, code, _, err := s.userDB.GetBudget(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "can't get user budget")
	}

	if budget == nil {
		return "У вас не задан бюджет на месяц", nil
	}

	value, err := s.converter.Exchange(ctx, *budget, converter.RUB, code)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Ваш бюджет на месяц: %.2f %s", value, code), nil
}

var budgetRe = regexp.MustCompile(`\/set_budget ([0-9.]+[0-9]+$)`)

// Set user budget. Will save it to database
func (s *Model) setBudget(ctx context.Context, msg Message) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "set budget")
	span.SetTag("command", commands.SetBudget)
	defer span.Finish()

	matches := budgetRe.FindStringSubmatch(msg.Text)
	if len(matches) < 2 {
		return ErrorIncorrectLine
	}

	budget, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return errors.Wrap(err, "can't conv budget to float")
	}

	code, err := s.getCode(ctx, msg.UserID)
	if err != nil {
		return errors.Wrap(err, "can't get code state")
	}

	budget, err = s.converter.Exchange(ctx, budget, code, converter.RUB)
	if err != nil {
		return errors.Wrap(err, "can't convert value")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = s.userDB.SetBudget(ctx, msg.UserID, budget)
	if err != nil {
		return errors.Wrap(err, "can't set budget")
	}

	return nil
}

var expenseRe = regexp.MustCompile(`^([0-9.]+) ([а-яА-Яa-zA-Z]+) ?([0-9]{4}-[0-9]{2}-[0-9]{2})?$`)

var ErrorIncorrectLine = errors.New("Incorrect line")

// Parse line with user expense
func parseExpense(ctx context.Context, text string) (*domain.Expense, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "parse user text")
	defer span.Finish()

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
