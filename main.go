package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/app"
	"github.com/taranovegor/naganbot/domain"
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

	err = a.ORM.AutoMigrate(
		&domain.Chat{},
		&domain.User{},
		&domain.Game{},
		&domain.Gunslinger{},
	)
	if err != nil {
		panic(err)
	}

	if a.ORM.Migrator().HasColumn(&domain.Chat{}, "required_players") {
		if err := a.ORM.Exec("UPDATE chats SET required_players = 6 WHERE required_players IS NULL").Error; err != nil {
			panic(err)
		}
		if err := a.ORM.Migrator().AlterColumn(&domain.Chat{}, "Settings.RequiredPlayers"); err != nil {
			panic(err)
		}
	}

	if a.ORM.Migrator().HasColumn(&domain.Game{}, "players_count") {
		if err := a.ORM.Exec("UPDATE games SET players_count = 6 WHERE players_count IS NULL").Error; err != nil {
			panic(err)
		}
		if err := a.ORM.Migrator().AlterColumn(&domain.Game{}, "PlayersCount"); err != nil {
			panic(err)
		}
	}

	log.Printf("authorized on account %s", a.BotAPI.Self.String())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := a.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		go safeExecute(func() {
			a.HandleUpdate(update)
		})
	}
}
