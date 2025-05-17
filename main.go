package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/container"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"gorm.io/gorm"
	"log"
)

var Version = "development"

func main() {
	fmt.Println(fmt.Sprintf("Nagan bot! Version: %s", Version))

	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load dotenv file")
	}

	sc, err := container.Init()
	if err != nil {
		panic(err)
	}

	orm := sc.Get(container.ORM).(*gorm.DB)
	err = orm.AutoMigrate(
		&domain.Chat{},
		&domain.User{},
		&domain.Game{},
		&domain.Gunslinger{},
	)
	if err != nil {
		panic(err)
	}

	trans := sc.Get(container.Translator).(*translator.Translator)
	chatRepository := sc.Get(container.RepositoryChat).(domain.ChatRepository)
	userRepository := sc.Get(container.RepositoryUser).(domain.UserRepository)
	bot := sc.Get(container.Bot).(*service.Bot)
	botApi := sc.Get(container.BotTelegram).(*tgbotapi.BotAPI)
	registry := sc.Get(container.CommandRegistry).(*command.Registry)

	log.Printf("authorized on account %s", botApi.Self.String())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message"}

	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		msg := update.Message
		if msg != nil {
			if !msg.IsCommand() {
				continue
			}

			chat := msg.Chat
			if chat.IsPrivate() || chat.IsChannel() {
				message := trans.Get("available only in chat", translator.Config{})
				bot.SendMessage(chat.ID, message)

				continue
			}

			domainChat := domain.NewChat(chat.ID, chat.Title, chat.UserName)
			if chatRepository.Exists(chat.ID) {
				chatRepository.Update(domainChat)
			} else {
				chatRepository.Store(domainChat)
			}

			from := msg.From
			domainUser := domain.NewUser(from.ID, from.FirstName, from.LastName, from.UserName)
			if userRepository.Exists(from.ID) {
				userRepository.Update(domainUser)
			} else {
				userRepository.Store(domainUser)
			}

			name := msg.Command()
			cmd, err := registry.Find(name)
			if err == nil {
				go cmd.Execute(msg)
			} else {
				log.Println(err.Error())
			}
		}
	}
}
