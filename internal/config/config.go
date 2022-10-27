package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type RateConfig struct {
	Key  string `yaml:"key"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Service struct {
	Tg   TelegramConfig `yaml:"telegram"`
	Rate RateConfig     `yaml:"rateApi"`
	DB   DatabaseConfig `yaml:"database"`
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

	err = yaml.Unmarshal(rawYAML, &s)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

// TELEGRAM

func (s *Service) Token() string {
	return s.Tg.Token
}

// RATE API

func (s *Service) Key() string {
	return s.Rate.Key
}

func (s *Service) Host() string {
	return s.Rate.Host
}

// DATABASE

func (s *Service) HostDB() string {
	return s.DB.Host
}

func (s *Service) Port() int {
	return s.DB.Port
}

func (s *Service) Username() string {
	return s.DB.Username
}

func (s *Service) Password() string {
	return s.DB.Password
}
