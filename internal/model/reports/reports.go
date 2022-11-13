package reports

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/commands"
	"go.uber.org/zap"
)

type ReportSender interface {
	SendReport(ctx context.Context, userID int64, report string) error
}

type ExpenseManipulator interface {
	Get(ctx context.Context, userID int64) ([]domain.Expense, error)
}

type UserManipulator interface {
	GetCode(ctx context.Context, userID int64) (string, error)
}

type ReportCacher interface {
	GetWeekReport(ctx context.Context, key int64) (string, bool)
	SetWeekReport(ctx context.Context, key int64, value string) error
	GetMonthReport(ctx context.Context, key int64) (string, bool)
	SetMonthReport(ctx context.Context, key int64, value string) error
	GetYearReport(ctx context.Context, key int64) (string, bool)
	SetYearReport(ctx context.Context, key int64, value string) error
}

type Converter interface {
	UpdateHistoricalRates(ctx context.Context, date *int64) error
	GetHistoricalCodeRate(ctx context.Context, code string, date int64) (float64, error)
}

type Model struct {
	grpcClient  ReportSender
	userDB      UserManipulator
	expenseDB   ExpenseManipulator
	reportCache ReportCacher
	converter   Converter
}

func New(grpcClient ReportSender, userDB UserManipulator, expenseDB ExpenseManipulator, reportCache ReportCacher, converter Converter) *Model {
	return &Model{
		grpcClient:  grpcClient,
		userDB:      userDB,
		expenseDB:   expenseDB,
		reportCache: reportCache,
		converter:   converter,
	}
}

func (s *Model) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Model) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Model) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		command := string(message.Key)
		logger.Info("get message from broker", zap.String("user_id", string(message.Value)), zap.String("command", command), zap.Int32("partition", message.Partition), zap.Int64("offset", message.Offset))

		userID, err := strconv.ParseInt(string(message.Value), 10, 64)
		if err != nil {
			logger.Error("cannot convert user id", zap.Error(err), zap.String("user_id", string(message.Value)))
		}

		report, err := s.getReportText(context.Background(), int64(userID), command)
		if err != nil {
			logger.Error("cannot calc user report", zap.Error(err), zap.Int64("user_id", userID), zap.String("command", command))
		}
		logger.Debug("report calculated", zap.Int64("user_id", userID), zap.String("report", report))

		err = s.grpcClient.SendReport(ctx, userID, report)
		if err != nil {
			logger.Error("cannot send user report", zap.Error(err), zap.Int64("user_id", userID), zap.String("command", command))
		}

		session.MarkMessage(message, "")
		cancel()
	}
	return nil
}

// Send prepared report with expenses to user
func (s *Model) getReportText(ctx context.Context, userID int64, text string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get report text")
	span.SetTag("command", commands.CommandReportText(text))
	defer span.Finish()

	currentTime := time.Now()
	var startTime time.Time
	var report string
	var ok bool
	switch text {
	case commands.CommandWeekReport:
		startTime = currentTime.AddDate(0, 0, -int(currentTime.Weekday())) // Start from Monday
		report, ok = s.reportCache.GetWeekReport(ctx, userID)
	case commands.CommandMonthReport:
		startTime = currentTime.AddDate(0, 0, 1-currentTime.Day()) // Start from first day in month
		report, ok = s.reportCache.GetMonthReport(ctx, userID)
	case commands.CommandYearReport:
		startTime = currentTime.AddDate(0, 1-int(currentTime.Month()), 1-currentTime.Day()) // Start with first dat in year
		report, ok = s.reportCache.GetYearReport(ctx, userID)
	}
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	if ok {
		return "Отчет:\n" + report, nil
	}

	report, err := s.calcReport(ctx, startTime, userID)
	if err != nil {
		return "", errors.Wrap(err, "can't get report")
	}

	if len(report) == 0 {
		return "Для начала добавьте покупки", nil
	}

	switch text {
	case commands.CommandWeekReport:
		err = s.reportCache.SetWeekReport(ctx, userID, report)
	case commands.CommandMonthReport:
		err = s.reportCache.SetMonthReport(ctx, userID, report)
	case commands.CommandYearReport:
		err = s.reportCache.SetYearReport(ctx, userID, report)
	}
	if err != nil {
		return "", err
	}
	return "Отчет:\n" + report, nil
}

// Get calculated user sum of expenses from `startTime`
func (s *Model) calcReport(ctx context.Context, startTime time.Time, userID int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "calc report")
	span.Finish()

	list, err := s.getExpenses(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "can't get expenses")
	}

	code, err := s.getCode(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "can't get current code")
	}

	sum := make(map[string]float64)
	for _, elem := range list {
		if elem.Date > startTime.Unix() {
			if code != converter.RUB {
				rate, err := s.converter.GetHistoricalCodeRate(ctx, code, elem.Date)
				switch {
				case err == sql.ErrNoRows:
					if err = s.converter.UpdateHistoricalRates(ctx, &elem.Date); err != nil {
						return "", err
					}

					rate, err = s.converter.GetHistoricalCodeRate(ctx, code, elem.Date)
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

// Get list of all user expenses
func (s *Model) getExpenses(ctx context.Context, userID int64) ([]domain.Expense, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get expenses")
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	list, err := s.expenseDB.Get(ctx, userID)
	return list, err
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
