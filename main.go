package main

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/app"
)

var Version = "development"

func safeExecute(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()
	fn()
}

func main() {
	fmt.Println(fmt.Sprintf("Nagan bot! Version: %s", Version))

	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load dotenv file")
	}

	a, err := app.New()
	if err != nil {
		panic(err)
	}

	if err := a.Migrate(); err != nil {
		panic(err)
	}

	log.Printf("authorized on account %s", a.BotAPI.Self.String())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := a.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		go safeExecute(func() {
			a.HandleUpdate(context.Background(), update)
		})
	}
}
