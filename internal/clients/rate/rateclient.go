package rate

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
)

type HeaderGetter interface {
	Key() string
	Host() string
}

type Service struct {
	Key  string
	Host string
}

func New(headerGetter HeaderGetter) *Service {
	return &Service{
		Key:  headerGetter.Key(),
		Host: headerGetter.Host(),
	}
}

func (s *Service) GetUpdate() (*converter.Rate, error) {
	rawJSON, err := s.getRequestRate()
	if err != nil {
		return nil, errors.Wrap(err, "can't complete get request")
	}

	responseRate, err := parseRates(rawJSON)
	if err != nil {
		return nil, errors.Wrap(err, "can't complete parse response")
	}

	rate := changeEURBaseToRUB(responseRate)

	return rate, nil
}

const url = "https://currency-conversion-and-exchange-rates.p.rapidapi.com/latest"

func (s *Service) getRequestRate() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get request exit with error")
	}

	req.Header.Add("X-RapidAPI-Key", s.Key)
	req.Header.Add("X-RapidAPI-Host", s.Host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "response error")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}

	return body, nil
}

type ResponseRate struct {
	Base  string `json:"base"` // EUR base
	Rates struct {
		USD float64 `json:"USD"`
		RUB float64 `json:"RUB"`
		CNY float64 `json:"CNY"`
	} `json:"rates"`
}

func parseRates(rawJSON []byte) (*ResponseRate, error) {
	responseRate := ResponseRate{}
	err := json.Unmarshal(rawJSON, &responseRate)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response json")
	}
	return &responseRate, nil
}

func changeEURBaseToRUB(responseRate *ResponseRate) *converter.Rate {
	return &converter.Rate{
		EUR: responseRate.Rates.RUB,
		USD: (1.0 / responseRate.Rates.USD) * responseRate.Rates.RUB,
		CNY: (1.0 / responseRate.Rates.CNY) * responseRate.Rates.RUB,
	}
}
