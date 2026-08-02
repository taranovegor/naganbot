package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/app"
)

var Version = "development"

func main() {
	fmt.Printf("Nagan bot! Version: %s\n", Version)

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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := a.BotAPI.GetUpdatesChan(u)

	stopReceivingUpdates := sync.OnceFunc(a.BotAPI.StopReceivingUpdates)
	go func() {
		<-ctx.Done()
		stopReceivingUpdates()
	}()

	a.Run(ctx, updates)
}
