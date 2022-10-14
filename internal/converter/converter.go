package converter

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type RateUpdater interface {
	GetUpdate(ctx context.Context) (*Rate, error)
}

type Rate struct {
	USD float64
	EUR float64
	CNY float64
}

type Model struct {
	rateClient RateUpdater
	rate       *Rate
}

func New(rateClient RateUpdater) *Model {
	return &Model{
		rateClient: rateClient,
		rate:       nil,
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
					break
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

	rate, err := m.rateClient.GetUpdate(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get update for rates")
	}
	m.setRate(rate)

	return nil
}

func (m *Model) setRate(rate *Rate) {
	m.rate = rate
}

func (m *Model) GetRate() Rate {
	return *m.rate
}

var CurrencyNotExistError = errors.New("currency not exist")

func (m *Model) getCurrencyRate(currency string) (float64, error) {
	if m.rate == nil {
		err := m.UpdateRate()
		if err != nil {
			return 0.0, errors.Wrap(err, "can't update rate")
		}
	}

	switch currency {
	case RUB:
		return 1.0, nil
	case USD:
		return m.rate.USD, nil
	case EUR:
		return m.rate.EUR, nil
	case CNY:
		return m.rate.CNY, nil
	default:
		return 0.0, CurrencyNotExistError
	}
}

func (m *Model) Exchange(amount float64, from string, to string) (float64, error) {
	fromRate, err := m.getCurrencyRate(from)
	if err != nil {
		return fromRate, errors.Wrap(err, "can't get from value in exchage")
	}

	toRate, err := m.getCurrencyRate(to)
	if err != nil {
		return toRate, errors.Wrap(err, "can't get to value in exchange")
	}

	return amount * (fromRate / toRate), nil
}
