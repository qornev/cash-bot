package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	RateApi struct {
		Key  string `yaml:"key"`
		Host string `yaml:"host"`
	} `yaml:"rateApi"`
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	return NewFromFile(configFile)
}

func NewFromFile(filePath string) (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Telegram.Token
}

func (s *Service) Key() string {
	return s.config.RateApi.Key
}

func (s *Service) Host() string {
	return s.config.RateApi.Host
}
