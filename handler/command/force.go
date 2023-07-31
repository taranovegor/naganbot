package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
)

type ForceHandler struct {
	Handler
	bot *service.Bot
}

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

func (hdlr ForceHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	hdlr.bot.Kick(chatID, msg.From.ID)
	hdlr.bot.DeleteMessage(chatID, msg.MessageID)
}
