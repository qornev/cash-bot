package main

import (
	"log"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/database"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/model/messages"
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

	db, err := database.New()
	if err != nil {
		log.Fatal("db init failed")
	}

	msgModel := messages.New(tgClient, db)

	tgClient.ListenUpdates(msgModel)
}
