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

type MetricsModel struct {
	server *http.Server
}

func NewMetricsServer(port int) *MetricsModel {
	model := &MetricsModel{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}

	http.Handle("/metrics", promhttp.Handler())
	return model
}

func (m *MetricsModel) StartMetricsServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting metrics http server...")
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("cannot start metrics http server", zap.Error(err))
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := m.server.Shutdown(ctx); err != nil {
			logger.Error("cannot shutdown metrics http server", zap.Error(err))
			return
		}
		logger.Info("http metrics server shutted down")
	}()
}
