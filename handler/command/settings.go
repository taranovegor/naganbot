package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type settingsHandler struct {
	chatRepo domain.ChatRepository
	trans    *translator.Translator
	bot      *service.Bot
}

func NewSettingsHandler(
	chatRepo domain.ChatRepository,
	trans *translator.Translator,
	bot *service.Bot,
) Handler {
	return &settingsHandler{
		chatRepo: chatRepo,
		trans:    trans,
		bot:      bot,
	}
}

func (hdlr *settingsHandler) Name() string {
	return "settings"
}

func (hdlr *settingsHandler) Execute(msg *tgbotapi.Message) {
	chat, err := hdlr.chatRepo.Get(msg.Chat.ID)
	if err != nil {
		hdlr.bot.SendMessage(msg.Chat.ID, hdlr.trans.Get("something went wrong", translator.Config{}))

		return
	}
	settings := chat.Settings

	message := hdlr.trans.Get("available settings below", translator.Config{})
	keyboard := callback.RevolverKeyboard(settings.RequiredPlayers, hdlr.trans)
	hdlr.bot.SendInlineKeyboard(msg.Chat.ID, message, keyboard)
}
