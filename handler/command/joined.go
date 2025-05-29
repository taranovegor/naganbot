package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/view"
)

type JoinedHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	game       domain.GameRepository
	msgFactory *view.JoinedMessageFactory
}

func NewJoinedHandler(
	bot *service.Bot,
	trans *translator.Translator,
	game domain.GameRepository,
	msgFactory *view.JoinedMessageFactory,
) Handler {
	return &JoinedHandler{
		bot:        bot,
		trans:      trans,
		game:       game,
		msgFactory: msgFactory,
	}
}

func (hdlr JoinedHandler) Name() string {
	return "joined"
}

func (hdlr JoinedHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	var message string
	var keyboard service.InlineKeyboard

	game, err := hdlr.game.GetLatestForChat(chatID)
	if err != nil {
		message = hdlr.trans.Get("active game not found", translator.Config{})
	} else {
		joinedMsg := hdlr.msgFactory.Create(game)
		message = joinedMsg.Text
		keyboard = joinedMsg.Keyboard
	}

	hdlr.bot.SendInlineKeyboard(chatID, message, keyboard)
}
