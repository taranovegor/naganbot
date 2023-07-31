package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const parseMode = tgbotapi.ModeHTML

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(
	api *tgbotapi.BotAPI,
) *Bot {
	return &Bot{
		api: api,
	}
}

func (bot Bot) SendMessage(chatID int64, text string) error {
	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat:  tgbotapi.BaseChat{ChatID: chatID},
		ParseMode: parseMode,
		Text:      text,
	})

	return err
}

func (bot Bot) DeleteMessage(chatID int64, messageID int) error {
	_, err := bot.api.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

	return err
}

func (bot Bot) Ban(chatID int64, userID int64, untilDate int64) error {
	_, err := bot.api.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: untilDate,
	})

	return err
}

func (bot Bot) Kick(chatID int64, userID int64) error {
	return bot.Ban(chatID, userID, time.Now().Add(time.Minute).Unix())
}
