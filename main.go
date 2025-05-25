package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/container"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
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

	if orm.Migrator().HasColumn(&domain.Chat{}, "required_players") {
		orm.Exec("UPDATE chats SET required_players = 6 WHERE required_players IS NULL")
		orm.Migrator().AlterColumn(&domain.Chat{}, "Settings.RequiredPlayers")
	}

	if orm.Migrator().HasColumn(&domain.Game{}, "players_count") {
		orm.Exec("UPDATE games SET players_count = 6 WHERE players_count IS NULL")
		orm.Migrator().AlterColumn(&domain.Game{}, "PlayersCount")
	}

	trans := sc.Get(container.Translator).(*translator.Translator)
	chatRepository := sc.Get(container.RepositoryChat).(domain.ChatRepository)
	userRepository := sc.Get(container.RepositoryUser).(domain.UserRepository)
	bot := sc.Get(container.Bot).(*service.Bot)
	botApi := sc.Get(container.BotTelegram).(*tgbotapi.BotAPI)
	cmdRegistry := sc.Get(container.CommandRegistry).(*command.Registry)
	clbRegistry := sc.Get(container.CallbackRegistry).(*callback.Registry)

	log.Printf("authorized on account %s", botApi.Self.String())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		chat := update.FromChat()
		if chat != nil {
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
		}

		from := update.SentFrom()
		if from != nil {
			domainUser := domain.NewUser(from.ID, from.FirstName, from.LastName, from.UserName)
			if userRepository.Exists(from.ID) {
				userRepository.Update(domainUser)
			} else {
				userRepository.Store(domainUser)
			}
		}

		if update.Message != nil {
			msg := update.Message
			if !msg.IsCommand() {
				continue
			}
			name := msg.Command()
			cmd, err := cmdRegistry.Find(name)
			if err == nil {
				go cmd.Execute(msg)
			} else {
				log.Println(err.Error())
			}
		} else if update.CallbackQuery != nil {
			callbackQuery := update.CallbackQuery
			query := callback.Pattern(callbackQuery.Data)
			hdlr, err := clbRegistry.Find(query)
			if err == nil {
				go hdlr.Execute(callbackQuery)
			} else {
				log.Println(err.Error())
			}
		}
	}
}
