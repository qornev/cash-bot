package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

type Model struct {
	server *http.Server
}

func NewServer(port int) *Model {
	model := &Model{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}

	http.Handle("/metrics", promhttp.Handler())
	return model
}

func (m *Model) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting http server...")
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("cannot start http server", zap.Error(err))
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := m.server.Shutdown(ctx); err != nil {
			logger.Error("cannot shutdown http server", zap.Error(err))
			return
		}
		logger.Info("http server shutted down")
	}()
}
