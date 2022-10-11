package rate

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
)

type ConfigGetter interface {
	Key() string
	Host() string
}

type Params struct {
	Key  string
	Host string
}

type Client struct {
	client *http.Client
	params Params
}

func New(configGetter ConfigGetter) *Client {
	return &Client{
		client: &http.Client{},
		params: Params{
			Key:  configGetter.Key(),
			Host: configGetter.Host(),
		},
	}
}

func (c *Client) GetUpdate(ctx context.Context) (*converter.Rate, error) {
	rawJSON, err := c.getRequestRate(ctx)
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

func (c *Client) getRequestRate(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get request exit with error")
	}

	req.Header.Add("X-RapidAPI-Key", c.params.Key)
	req.Header.Add("X-RapidAPI-Host", c.params.Host)

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
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
