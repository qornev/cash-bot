package middlewares

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
)

func InitTracing(serviceName string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	return err
}

func TracingMiddleware(next MessageProcesser) MessageProcesser {
	return MessageProcesserFunc(func(msg messages.Message, info *messages.CommandInfo) error {
		ctx := info.Context()

		span, ctx := opentracing.StartSpanFromContext(ctx, "incoming message")
		defer span.Finish()

		info = info.WithContext(ctx)
		err := next.IncomingMessage(msg, info)
		return err
	})
}
