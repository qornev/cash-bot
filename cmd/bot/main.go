package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/rate"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed")
	}

	storage, err := storage.New()
	if err != nil {
		log.Fatal("storage init failed")
	}

	rateClient := rate.New(config)
	converter := converter.New(rateClient)

	msgModel := messages.New(tgClient, storage, converter)
	clbModel := callbacks.New(tgClient, storage)

	wg := sync.WaitGroup{}
	converter.AutoUpdateRate(ctx, &wg)
	tgClient.AutoListenUpdates(ctx, &wg, msgModel, clbModel)

	<-ctx.Done()
	wg.Wait()
	log.Println("all process are finished")
}
