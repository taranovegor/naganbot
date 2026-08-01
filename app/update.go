package app

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/translator"
)

func (a *App) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	chat := update.FromChat()
	if chat != nil {
		if chat.IsPrivate() || chat.IsChannel() {
			message := a.Translator.Get("available only in chat", translator.Config{})
			a.Bot.SendMessage(chat.ID, message)

			return
		}

		if err := a.Chats.Save(domain.NewChat(chat.ID, chat.Title, chat.UserName)); err != nil {
			log.Printf("failed to save chat %d: %v", chat.ID, err)
		}
	}

	from := update.SentFrom()
	if from != nil {
		if err := a.Users.Save(domain.NewUser(from.ID, from.FirstName, from.LastName, from.UserName)); err != nil {
			log.Printf("failed to save user %d: %v", from.ID, err)
		}
	}

	if update.Message != nil {
		msg := update.Message
		if !msg.IsCommand() {
			return
		}
		name := msg.Command()
		cmd, err := a.Commands.Find(name)
		if err != nil {
			log.Println(err.Error())
			return
		}
		cmd.Execute(ctx, msg)
	} else if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		query := callback.Pattern(callbackQuery.Data)
		hdlr, err := a.Callbacks.Find(query)
		if err != nil {
			log.Println(err.Error())
			return
		}
		hdlr.Execute(ctx, callbackQuery)
	}
}
