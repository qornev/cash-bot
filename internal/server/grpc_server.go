package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/gateway/v1"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageSender interface {
	SendMessage(ctx context.Context, text string, userID int64) error
}

type GRPCModel struct {
	client MessageSender
	server *grpc.Server
	lis    net.Listener
	gateway.UnimplementedRouterServer
}

func NewGRPCServer(port int, client MessageSender) (*GRPCModel, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	model := &GRPCModel{
		client: client,
		server: s,
		lis:    lis,
	}
	gateway.RegisterRouterServer(s, model)
	return model, nil
}

func (s *GRPCModel) SendReport(ctx context.Context, report *gateway.Report) (*empty.Empty, error) {
	err := s.client.SendMessage(ctx, report.Text, report.UserId)
	return &emptypb.Empty{}, err
}

func (s *GRPCModel) StartGRPCServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting grpc server...")
		reflection.Register(s.server)
		if err := s.server.Serve(s.lis); err != nil {
			logger.Fatal("cannot start grpc server", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.server.GracefulStop()
		logger.Info("grpc server shutted down")
	}()
}
