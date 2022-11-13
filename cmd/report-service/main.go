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
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	flag.Parse()

	if err := logger.InitLogger(*developMode); err != nil {
		log.Fatal("logger init failed", err.Error())
	}

	config, err := config.New()
	if err != nil {
		logger.Fatal("config init failed", zap.Error(err))
	}

	rateClient := rate.New(config)

	db, err := storage.Connect(config)
	if err != nil {
		logger.Fatal("db client init failed", zap.Error(err))
	}
	userDB := storage.NewUserDB(db)
	expenseDB := storage.NewExpenseDB(db)
	rateDB := storage.NewRateDB(db)

	reportCache := cache.New(config)

	converter := converter.New(rateClient, rateDB, userDB, reportCache)

	reportModel := reports.New(userDB, expenseDB, reportCache, converter)

	consumer, err := consumer.NewConsumer(config.ListBroker(), title.ReportsBuilder)
	if err != nil {
		logger.Error("consumer init failed", zap.Error(err))
	}

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
