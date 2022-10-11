package converter

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
)

type RateUpdater interface {
	GetUpdate() (*Rate, error)
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

func (m *Model) AutoUpdateRate() {
	ticker := time.NewTicker(time.Hour)

	go func() {
		for {
			<-ticker.C

			if err := m.UpdateRate(); err != nil {
				log.Println("error processing rate update:", err)
			}
		}
	}()
}

func (m *Model) UpdateRate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rate := make(chan *Rate, 1)  // Add buffer with size 1 to prevent
	chErr := make(chan error, 1) // infinity goroutine in line 51

	go func() {
		value, err := m.rateClient.GetUpdate()
		if err != nil {
			chErr <- err
			return
		}
		rate <- value
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case value := <-rate:
		m.setRate(value)
	case errValue := <-chErr:
		return errValue
	}

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
