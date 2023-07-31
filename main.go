package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/handler/command"
	"log"
	"os"
)

const Version = "2.1.0"

func main() {
	fmt.Println(fmt.Sprintf("Nagan bot! Version: %s", Version))

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("APP_TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	log.Printf("authorized on account %s", bot.Self.String())

	registry := command.NewRegistry(
		"ng",
		command.NewForce(bot),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message"}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			name := update.Message.Command()
			cmd, err := registry.Find(name)
			if nil == err {
				err = cmd.Execute(update.Message)
			}

			if nil != err {
				log.Printf("[%s] command execution result: %s", name, err.Error())
			}
		}
	}
}
