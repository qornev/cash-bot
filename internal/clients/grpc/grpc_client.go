package grpc

import (
	"context"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/gateway/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	conn   *grpc.ClientConn
	client gateway.RouterClient
}

func New(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := gateway.NewRouterClient(conn)
	return &Service{
		conn:   conn,
		client: client,
	}, nil
}

func (s *Service) SendReport(ctx context.Context, userID int64, report string) error {
	_, err := s.client.SendReport(ctx, &gateway.Report{
		UserId: userID,
		Text:   report,
	})
	return err
}

func (s *Service) Close() {
	s.conn.Close()
}
