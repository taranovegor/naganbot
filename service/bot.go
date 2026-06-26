package service

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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

func (bot Bot) SendMessage(chatID int64, text string) {
	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: chatID},
		ParseMode:             parseMode,
		Text:                  text,
		DisableWebPagePreview: true,
	})
	if err != nil {
		log.Printf("failed to send message to chat %d: %v", chatID, err)
	}
}

func (bot Bot) DeleteMessage(chatID int64, messageID int) {
	_, err := bot.api.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
	if err != nil {
		log.Printf("failed to delete message %d in chat %d: %v", messageID, chatID, err)
	}
}

func (bot Bot) Ban(chatID int64, userID int64, untilDate int64) error {
	_, err := bot.api.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: untilDate,
	})
	if err != nil {
		log.Printf("failed to ban user %d in chat %d: %v", userID, chatID, err)
	}

	return err
}

func (bot Bot) Kick(chatID int64, userID int64) error {
	return bot.Ban(chatID, userID, time.Now().Add(time.Minute).Unix())
}

func (bot Bot) SendInlineKeyboard(chatID int64, text string, keyboard []map[string]string) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, row := range keyboard {
		var cols []tgbotapi.InlineKeyboardButton
		for key, val := range row {
			cols = append(cols, tgbotapi.NewInlineKeyboardButtonData(val, key))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cols...))
	}

	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(rows...),
		},
		ParseMode: parseMode,
		Text:      text,
	})
	if err != nil {
		log.Printf("failed to send inline keyboard to chat %d: %v", chatID, err)
	}
}

func (bot Bot) EditMessageReplyMarkup(chatID int64, messageID int, keyboard []map[string]string) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, row := range keyboard {
		var cols []tgbotapi.InlineKeyboardButton
		for key, val := range row {
			cols = append(cols, tgbotapi.NewInlineKeyboardButtonData(val, key))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cols...))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	_, err := bot.api.Request(tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, markup))
	if err != nil {
		log.Printf("failed to edit message %d reply markup in chat %d: %v", messageID, chatID, err)
	}
}

func (bot Bot) AnswerCallback(callbackQueryID string, text string) {
	_, err := bot.api.Request(tgbotapi.NewCallback(callbackQueryID, text))
	if err != nil {
		log.Printf("failed to answer callback %s: %v", callbackQueryID, err)
	}
}

func (bot Bot) IsAdmin(chatID int64, userID int64) (bool, error) {
	resp, err := bot.api.Request(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	})
	if err != nil {
		return false, err
	}

	var records []struct {
		User struct {
			ID int64 `json:"id"`
		} `json:"user"`
	}

	err = json.Unmarshal(resp.Result, &records)
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record.User.ID == userID {
			return true, nil
		}
	}

	return false, nil
}
