package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/cache"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/grpc"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/rate"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/kafka/consumer"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/kafka/title"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/reports"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/storage"
	"go.uber.org/zap"
)

var (
	developMode = flag.Bool("develop", false, "development mode")
	addr        = flag.String("addr", "127.0.0.1:8081", "address of grpc server")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	flag.Parse()

	// INIT LOGGER
	if err := logger.InitLogger(*developMode); err != nil {
		log.Fatal("logger init failed", err.Error())
	}

	// INIT CONFIGS
	config, err := config.New()
	if err != nil {
		logger.Fatal("config init failed", zap.Error(err))
	}

	// INIT CLIENTS
	rateClient := rate.New(config)
	grpcClient, err := grpc.New(*addr)
	if err != nil {
		logger.Fatal("grpc client init failed", zap.Error(err))
	}
	defer grpcClient.Close()

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

	// INIT CONVERTER
	converter := converter.New(rateClient, rateDB, userDB, reportCache)

	// INIT MODELS
	reportModel := reports.New(grpcClient, userDB, expenseDB, reportCache, converter)

	// INIT CONSUMER
	consumer, err := consumer.NewConsumer(config.ListBroker(), title.ReportsBuilder)
	if err != nil {
		logger.Error("consumer init failed", zap.Error(err))
	}

	// START WORKERS
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start consume with group")
		if err := consumer.StartConsume(ctx, title.Reports, reportModel); err != nil {
			logger.Error("consume start error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	wg.Wait()
	logger.Info("cosumer shutted down gracefully")
}
