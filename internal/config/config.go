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

type CacheConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Service struct {
	Tg     TelegramConfig `yaml:"telegram"`
	Rate   RateConfig     `yaml:"rateApi"`
	DB     DatabaseConfig `yaml:"database"`
	Cache  CacheConfig    `yaml:"cache"`
	Broker []string       `yaml:"broker"`
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

func (s *Service) TokenTG() string {
	return s.Tg.Token
}

// RATE API

func (s *Service) KeyRateAPI() string {
	return s.Rate.Key
}

func (s *Service) HostRateAPI() string {
	return s.Rate.Host
}

// DATABASE

func (s *Service) HostDB() string {
	return s.DB.Host
}

func (s *Service) PortDB() int {
	return s.DB.Port
}

func (s *Service) UsernameDB() string {
	return s.DB.Username
}

func (s *Service) PasswordDB() string {
	return s.DB.Password
}

// CACHE

func (s *Service) HostCache() string {
	return s.Cache.Host
}

func (s *Service) PortCache() int {
	return s.Cache.Port
}

func (s *Service) UsernameCache() string {
	return s.Cache.Username
}

func (s *Service) PasswordCache() string {
	return s.Cache.Password
}

// Broker

func (s *Service) ListBroker() []string {
	return s.Broker
}
