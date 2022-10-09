package converter

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type HeaderGetter interface {
	Key() string
	Host() string
}

type Rate struct {
	USD float64
	EUR float64
	CNY float64
}

type Params struct {
	Key  string
	Host string
}

type Service struct {
	rate   *Rate
	params Params
}

func New(headerGetter HeaderGetter) *Service {
	return &Service{
		rate: nil,
		params: Params{
			Key:  headerGetter.Key(),
			Host: headerGetter.Host(),
		},
	}
}

func (s *Service) UpdateRate() error {
	rawJSON, err := s.getRequestRate()
	if err != nil {
		return errors.Wrap(err, "can't complete get request")
	}

	responseRate, err := parseRates(rawJSON)
	if err != nil {
		return errors.Wrap(err, "can't complete parse response")
	}

	rate := changeEURBaseToRUB(responseRate)
	s.setRate(rate)

	return nil
}

const url = "https://currency-conversion-and-exchange-rates.p.rapidapi.com/latest"

func (s *Service) getRequestRate() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get request exit with error")
	}

	req.Header.Add("X-RapidAPI-Key", s.params.Key)
	req.Header.Add("X-RapidAPI-Host", s.params.Host)

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

func (s *Service) setRate(rate Rate) {
	s.rate = &rate
}

var CurrencyNotExistError = errors.New("currency not exist")

func (s *Service) getRate(currency string) (float64, error) {
	switch currency {
	case "RUB":
		return 1.0, nil
	case "USD":
		return s.rate.USD, nil
	case "EUR":
		return s.rate.EUR, nil
	case "CNY":
		return s.rate.CNY, nil
	default:
		return 0.0, CurrencyNotExistError
	}
}

func (s *Service) Exchange(value float64, from string, to string) (float64, error) {
	fromRate, err := s.getRate(from)
	if err != nil {
		return fromRate, errors.Wrap(err, "can't get from value in exchage")
	}

	toRate, err := s.getRate(to)
	if err != nil {
		return toRate, errors.Wrap(err, "can't get to value in exchange")
	}

	return value * (fromRate / toRate), nil
}

type ResponseRate struct {
	EUR   string `json:"base"`
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

func changeEURBaseToRUB(responseRate *ResponseRate) Rate {
	return Rate{
		EUR: responseRate.Rates.RUB,
		USD: (1.0 / responseRate.Rates.USD) * responseRate.Rates.RUB,
		CNY: (1.0 / responseRate.Rates.CNY) * responseRate.Rates.RUB,
	}
}
