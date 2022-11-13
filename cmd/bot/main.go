package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/cache"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/rate"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/kafka/producer"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/middlewares"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/server"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/storage"
	"go.uber.org/zap"
)

var (
	metricsPort = flag.Int("mport", 8080, "port to listen metrics")
	grpcPort    = flag.Int("gport", 8081, "port to listen grpc requests")
	developMode = flag.Bool("develop", false, "development mode")
	serviceName = flag.String("service", "tgbot", "name of starting service")
)

func main() {
	// Initialize context with kill and interruption processes
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	flag.Parse()

	// INIT LOGGER
	if err := logger.InitLogger(*developMode); err != nil {
		log.Fatal("logger init failed", err.Error())
	}

	// INIT TRACING
	if err := middlewares.InitTracing(*serviceName); err != nil {
		logger.Fatal("tracing init failed", zap.Error(err))
	}

	// INIT CONFIG
	config, err := config.New()
	if err != nil {
		logger.Fatal("config init failed:", zap.Error(err))
	}

	// INIT CLIENTS
	tgClient, err := tg.New(config)
	if err != nil {
		logger.Fatal("tg client init failed", zap.Error(err))
	}

	rateClient := rate.New(config)

	// INIT DATABASES
	db, err := storage.Connect(config)
	if err != nil {
		logger.Fatal("db client init failed", zap.Error(err))
	}
	userDB := storage.NewUserDB(db)
	expenseDB := storage.NewExpenseDB(db)
	rateDB := storage.NewRateDB(db)

	// INIT CACHE
	reportCache := cache.New(config)

	// INIT PRODUCER
	producer, err := producer.NewProducer(config.ListBroker())
	if err != nil {
		logger.Fatal("message producer init failed", zap.Error(err))
	}

	// INIT MODELS
	converter := converter.New(rateClient, rateDB, userDB, reportCache)

	msgModel := messages.New(tgClient, userDB, expenseDB, reportCache, converter, producer)
	clbModel := callbacks.New(tgClient, userDB, reportCache)

	// INIT SERVER
	metricsServer := server.NewMetricsServer(*metricsPort)
	GRPCServer, err := server.NewGRPCServer(*grpcPort)
	if err != nil {
		logger.Fatal("grpc server init failed", zap.Error(err))
	}

	// START WORKERS
	wg := sync.WaitGroup{}
	converter.AutoUpdateRate(ctx, &wg)
	tgClient.AutoListenUpdates(ctx, &wg, msgModel, clbModel)
	metricsServer.StartMetricsServer(ctx, &wg)
	GRPCServer.StartGRPCServer(ctx, &wg)

	<-ctx.Done()
	wg.Wait()
	logger.Info("all processes are finished gracefully")
}
