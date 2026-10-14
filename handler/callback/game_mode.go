package callback

import (
	"context"
	"fmt"
	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/translator"
)

type gameMode struct {
	chatRepo domain.ChatRepository
	bot      messenger
	trans    *translator.Translator
}

var _ Handler = (*gameMode)(nil)

func NewGameMode(
	chatRepo domain.ChatRepository,
	bot messenger,
	trans *translator.Translator,
) Handler {
	return &gameMode{
		chatRepo: chatRepo,
		bot:      bot,
		trans:    trans,
	}
}

func (hdlr *gameMode) Pattern() Pattern {
	return GameMode
}

func (hdlr *gameMode) Execute(_ context.Context, query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	isAdmin, err := hdlr.bot.IsAdmin(chatID, query.From.ID)
	if err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if !isAdmin {
		notification := hdlr.trans.Get("settings can be changed only by admins", translator.Config{})
		hdlr.bot.AnswerCallback(query.ID, notification)
		return
	}

	mode := domain.GameMode(GameMode.GetArg(query.Data, 1))
	if !slices.Contains(gameModes, mode) {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	chat, err := hdlr.chatRepo.Get(chatID)
	if err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if chat.Settings.Mode == mode {
		hdlr.bot.AnswerCallback(query.ID, "")
		return
	}

	chat.Settings.Mode = mode
	if err := hdlr.chatRepo.UpdateSettings(&chat); err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	notification := fmt.Sprintf(
		"%s\n%s",
		hdlr.trans.Get(fmt.Sprintf("%s mode has been set", mode), translator.Config{}),
		hdlr.trans.Get("settings will be applied for next games", translator.Config{}),
	)
	hdlr.bot.AnswerCallback(query.ID, notification)

	keyboard := SettingsKeyboard(mode, chat.Settings.RequiredPlayers, hdlr.trans)
	hdlr.bot.EditMessageReplyMarkup(chatID, query.Message.MessageID, keyboard)
}
