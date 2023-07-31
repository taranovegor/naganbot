package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type CreateHandler struct {
	Handler
	bot   *service.Bot
	trans *translator.Translator
	base  *JoinHandler
}

func NewCreateHandler(
	bot *service.Bot,
	trans *translator.Translator,
	base *JoinHandler,
) Handler {
	return &CreateHandler{
		bot:   bot,
		trans: trans,
		base:  base,
	}
}

func (hdlr CreateHandler) Name() string {
	return "create"
}

func (hdlr CreateHandler) Execute(msg *tgbotapi.Message) {
	hdlr.bot.SendMessage(msg.Chat.ID, hdlr.trans.Get("deprecated command", translator.Config{
		Args: map[string]string{
			"%new": "/ngjoin@naganbot",
		},
	}))

	hdlr.base.Execute(msg)
}
