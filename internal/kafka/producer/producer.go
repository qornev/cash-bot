package producer

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

type Service struct {
	producer sarama.SyncProducer
}

func NewProducer(brokerList []string) (*Service, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = time.Millisecond * 250
	config.Producer.Return.Successes = true
	if config.Producer.Idempotent {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &Service{
		producer: producer,
	}, nil
}

func (s *Service) ProduceMessage(topic string, userID int64, text string) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(text),
		Value: sarama.StringEncoder(fmt.Sprint(userID)),
	}

	partition, offset, err := s.producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	logger.Info(
		"Write message to broker",
		zap.String("topic", topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
	return nil
}

func (s *Service) Close() error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	return nil
}
