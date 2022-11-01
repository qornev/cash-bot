package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/rate"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/middlewares"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/server"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/storage"
	"go.uber.org/zap"
)

var (
	port        = flag.Int("port", 8080, "the port to listen")
	developMode = flag.Bool("develop", false, "development mode")
	serviceName = flag.String("service", "tgbot", "the name of starting service")
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

	// INIT MODELS
	converter := converter.New(rateClient, rateDB, userDB)

	msgModel := messages.New(tgClient, userDB, expenseDB, converter)
	clbModel := callbacks.New(tgClient, userDB)

	// INIT SERVER
	srv := server.NewServer(*port)

	// START WORKERS
	wg := sync.WaitGroup{}
	converter.AutoUpdateRate(ctx, &wg)
	tgClient.AutoListenUpdates(ctx, &wg, msgModel, clbModel)
	srv.Start(ctx, &wg)

	<-ctx.Done()
	wg.Wait()
	logger.Info("all processes are finished gracefully")
}
