package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
)

type ForceHandler struct {
	bot *service.Bot
}

var _ Handler = (*ForceHandler)(nil)

func NewForceHandler(
	bot *service.Bot,
) Handler {
	return &ForceHandler{
		bot: bot,
	}
}

func (hdlr ForceHandler) Name() string {
	return "force"
}

func (hdlr ForceHandler) Execute(_ context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	_ = hdlr.bot.Kick(chatID, msg.From.ID)
	hdlr.bot.DeleteMessage(chatID, msg.MessageID)
}
