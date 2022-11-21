package consumer

import (
	"context"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
)

type Service struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(brokerList []string, groupID string) (*Service, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Service{
		consumerGroup: consumerGroup,
	}, nil
}

func (s *Service) StartConsume(ctx context.Context, topic string, model sarama.ConsumerGroupHandler) error {
	// `Consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	defer s.consumerGroup.Close()
	for {
		logger.Info("join to cluster of consumers")
		err := s.consumerGroup.Consume(ctx, []string{topic}, model)
		switch {
		case err != nil:
			return err
		case ctx.Err() == context.Canceled:
			return nil
		}
	}
}
