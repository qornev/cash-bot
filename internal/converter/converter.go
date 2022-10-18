package converter

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type RateUpdater interface {
	GetUpdate(ctx context.Context) (*CurrentRate, error)
}

type RateManipulator interface {
	Add(ctx context.Context, date int64, code string, nominal float64) error
}

type CurrentRate struct {
	USD float64
	EUR float64
	CNY float64
}

type Rate struct {
	Code    string
	Nominal float64
}

type Model struct {
	rateClient  RateUpdater
	rateDB      RateManipulator
	currentRate *CurrentRate
}

func New(rateClient RateUpdater, rateDB RateManipulator) *Model {
	return &Model{
		rateClient:  rateClient,
		rateDB:      rateDB,
		currentRate: nil,
	}
}

func (m *Model) AutoUpdateRate(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Hour)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				if err := m.UpdateRate(); err != nil {
					log.Println("error processing rate update:", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *Model) UpdateRate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	currentRate, err := m.rateClient.GetUpdate(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get update for rates")
	}

	if err := m.setCurrentRate(ctx, currentRate); err != nil {
		return errors.Wrap(err, "can't set current rate")
	}

	return nil
}

func (m *Model) setCurrentRate(ctx context.Context, currentRate *CurrentRate) error {
	currentTime := time.Now().Unix()
	if err := m.rateDB.Add(ctx, currentTime, USD, currentRate.USD); err != nil {
		return err
	}
	if err := m.rateDB.Add(ctx, currentTime, EUR, currentRate.EUR); err != nil {
		return err
	}
	if err := m.rateDB.Add(ctx, currentTime, CNY, currentRate.CNY); err != nil {
		return err
	}

	m.currentRate = currentRate
	return nil
}

func (m *Model) GetCurrentRate() CurrentRate {
	return *m.currentRate
}

var CodeNotExistError = errors.New("code not exist")

func (m *Model) getCodeRate(code string) (float64, error) {
	if m.currentRate == nil {
		err := m.UpdateRate()
		if err != nil {
			return 0.0, errors.Wrap(err, "can't update rate")
		}
	}

	switch code {
	case RUB:
		return 1.0, nil
	case USD:
		return m.currentRate.USD, nil
	case EUR:
		return m.currentRate.EUR, nil
	case CNY:
		return m.currentRate.CNY, nil
	default:
		return 0.0, CodeNotExistError
	}
}

func (m *Model) Exchange(amount float64, from string, to string) (float64, error) {
	fromRate, err := m.getCodeRate(from)
	if err != nil {
		return fromRate, errors.Wrap(err, "can't get from value in exchage")
	}

	toRate, err := m.getCodeRate(to)
	if err != nil {
		return toRate, errors.Wrap(err, "can't get to value in exchange")
	}

	return amount * (fromRate / toRate), nil
}
