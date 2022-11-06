package converter

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

type RateUpdater interface {
	GetUpdate(ctx context.Context, date *int64) (*Rates, error)
}

type RateManipulator interface {
	Add(ctx context.Context, date int64, code string, nominal float64) error
	Get(ctx context.Context, date int64, code string) (*domain.Rate, error)
}

type UserManipulator interface {
	UpdateBudget(ctx context.Context, userID int64, budget float64) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type ReportCacher interface {
	RemoveFromAll(ctx context.Context, key []int64) error
}

type Rates struct {
	USD float64
	EUR float64
	CNY float64
}

type Model struct {
	rateClient   RateUpdater
	rateDB       RateManipulator
	userDB       UserManipulator
	reportCache  ReportCacher
	currentRates *Rates
}

func New(rateClient RateUpdater, rateDB RateManipulator, userDB UserManipulator, reportCache ReportCacher) *Model {
	return &Model{
		rateClient:   rateClient,
		rateDB:       rateDB,
		userDB:       userDB,
		reportCache:  reportCache,
		currentRates: nil,
	}
}

// Setup worker for auto updating
func (m *Model) AutoUpdateRate(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Hour)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				if err := m.UpdateRecentRates(ctx); err != nil {
					logger.Error("error processing rate update", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Get actual rates update from API
func (m *Model) UpdateRecentRates(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update recent rates")
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rates, err := m.rateClient.GetUpdate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "can't get update for rates")
	}

	if err := m.setCurrentRates(ctx, rates); err != nil {
		return errors.Wrap(err, "can't set current rates")
	}

	return nil
}

// Get rates update from API at specified `date`
func (m *Model) UpdateHistoricalRates(ctx context.Context, date *int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update historical rates")
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rates, err := m.rateClient.GetUpdate(ctx, date)
	if err != nil {
		return errors.Wrap(err, "can't get update for rates")
	}

	if err := m.addRates(ctx, rates, date); err != nil {
		return errors.Wrap(err, "can't add rates")
	}

	return nil
}

// Set rates in cache, also with add to database and updating budgets with not RUB currencies
func (m *Model) setCurrentRates(ctx context.Context, rates *Rates) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "set current rates")
	defer span.Finish()

	if err := m.addRates(ctx, rates, nil); err != nil {
		return err
	}

	users, err := m.userDB.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Budget == nil || user.Code == RUB {
			continue
		}

		rate, err := m.GetHistoricalCodeRate(ctx, user.Code, user.Updated)
		switch {
		case err == sql.ErrNoRows:
			if err = m.UpdateHistoricalRates(ctx, &user.Updated); err != nil {
				return err
			}

			rate, err = m.GetHistoricalCodeRate(ctx, user.Code, user.Updated)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		}

		var diff float64
		switch user.Code {
		case USD:
			diff = rates.USD / rate
		case EUR:
			diff = rates.EUR / rate
		case CNY:
			diff = rates.CNY / rate
		default:
			return CodeNotExistError
		}

		if err = m.userDB.UpdateBudget(ctx, user.ID, *user.Budget*diff); err != nil {
			return err
		}
	}

	m.currentRates = rates
	return nil
}

// Add rates to database
func (m *Model) addRates(ctx context.Context, rates *Rates, date *int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add rates")
	defer span.Finish()

	if date == nil {
		currentDate := time.Now().Unix()
		date = &currentDate
	}
	if err := m.rateDB.Add(ctx, *date, USD, rates.USD); err != nil {
		return err
	}
	if err := m.rateDB.Add(ctx, *date, EUR, rates.EUR); err != nil {
		return err
	}
	if err := m.rateDB.Add(ctx, *date, CNY, rates.CNY); err != nil {
		return err
	}
	return nil
}

// Get rates from cache
func (m *Model) GetCurrentRates() Rates {
	return *m.currentRates
}

var CodeNotExistError = errors.New("code not exist")

// Get currency rate with `code` from cache
func (m *Model) getCurrentCodeRate(ctx context.Context, code string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get current rate for code")
	span.SetTag("code", code)
	defer span.Finish()

	if m.currentRates == nil {
		if err := m.UpdateRecentRates(ctx); err != nil {
			return 0.0, errors.Wrap(err, "can't update rate")
		}
	}

	switch code {
	case RUB:
		return 1.0, nil
	case USD:
		return m.currentRates.USD, nil
	case EUR:
		return m.currentRates.EUR, nil
	case CNY:
		return m.currentRates.CNY, nil
	default:
		return 0.0, CodeNotExistError
	}
}

// Get currency rate with `code` from database
func (m *Model) GetHistoricalCodeRate(ctx context.Context, code string, date int64) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get historical rate for code")
	span.SetTag("code", code)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rate, err := m.rateDB.Get(ctx, date, code)
	if err != nil {
		return 0.0, err
	}

	return rate.Nominal, nil
}

// Exchage money with rates from cache. If no rates in cache, method will request it
func (m *Model) Exchange(ctx context.Context, amount float64, from string, to string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "exchange money")
	defer span.Finish()

	fromRate, err := m.getCurrentCodeRate(ctx, from)
	if err != nil {
		return fromRate, errors.Wrap(err, "can't get from value in exchage")
	}

	toRate, err := m.getCurrentCodeRate(ctx, to)
	if err != nil {
		return toRate, errors.Wrap(err, "can't get to value in exchange")
	}

	return amount * (fromRate / toRate), nil
}
