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

type Service struct {
	tg   TelegramConfig `yaml:"telegram"`
	rate RateConfig     `yaml:"rateApi"`
	db   DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
	return s.tg.Token
}

// RATE API

func (s *Service) Key() string {
	return s.rate.Key
}

func (s *Service) Host() string {
	return s.rate.Host
}

// DATABASE

func (s *Service) HostDB() string {
	return s.db.Host
}

func (s *Service) Port() int {
	return s.db.Port
}

func (s *Service) Username() string {
	return s.db.Username
}

func (s *Service) Password() string {
	return s.db.Password
}
