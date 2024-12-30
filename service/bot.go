package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const parseMode = tgbotapi.ModeHTML

type Bot interface {
	SendMessage(chatID int64, text string) error
	DeleteMessage(chatID int64, messageID int) error
	Ban(chatID int64, userID int64, untilDate int64) error
	Kick(chatID int64, userID int64) error
}

type bot struct {
	Bot
	api *tgbotapi.BotAPI
}

func NewBot(
	api *tgbotapi.BotAPI,
) Bot {
	return bot{
		api: api,
	}
}

func (bot bot) SendMessage(chatID int64, text string) error {
	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat:  tgbotapi.BaseChat{ChatID: chatID},
		ParseMode: parseMode,
		Text:      text,
	})

	return err
}

func (bot bot) DeleteMessage(chatID int64, messageID int) error {
	_, err := bot.api.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

	return err
}

func (bot bot) Ban(chatID int64, userID int64, untilDate int64) error {
	_, err := bot.api.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: untilDate,
	})

	return err
}

func (bot bot) Kick(chatID int64, userID int64) error {
	return bot.Ban(chatID, userID, time.Now().Add(time.Minute).Unix())
}
