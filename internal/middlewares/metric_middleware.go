package middlewares

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

var (
	HistogramResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tgbot",
			Subsystem: "command",
			Name:      "histogram_process_time_seconds",
			Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
		},
		[]string{"command"},
	)
)

func MetricMiddleware(next MessageProcesser) MessageProcesser {
	return MessageProcesserFunc(func(msg messages.Message, info *messages.CommandInfo) error {
		startTime := time.Now()
		err := next.IncomingMessage(msg, info)
		duration := time.Since(startTime)

		HistogramResponseTime.
			WithLabelValues(info.Command).
			Observe(duration.Seconds())
		return err
	})
}
