package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
)

type forceHandler struct {
	bot *service.Bot
}

var _ Handler = (*forceHandler)(nil)

func NewForceHandler(
	bot *service.Bot,
) Handler {
	return &forceHandler{
		bot: bot,
	}
}

func (hdlr *forceHandler) Name() string {
	return "force"
}

func (hdlr *forceHandler) Execute(_ context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	_ = hdlr.bot.Kick(chatID, msg.From.ID)
	hdlr.bot.DeleteMessage(chatID, msg.MessageID)
}
