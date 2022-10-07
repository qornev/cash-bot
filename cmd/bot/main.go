package main

import (
	"log"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/callbacks"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/storage"
)

func main() {
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

	msgModel := messages.New(tgClient, storage)
	clbModel := callbacks.New(tgClient, storage)

	tgClient.ListenUpdates(msgModel, clbModel)
}
