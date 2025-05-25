package command

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
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
		hdlr.bot.SendMessage(chat.ID, hdlr.trans.Get("something went wrong", translator.Config{}))

		return
	}
	settings := chat.Settings

	message := hdlr.trans.Get("available settings below", translator.Config{})
	var keyboard []map[string]string
	for _, shotNum := range []int{4, 6, 7} {
		shotStr := strconv.Itoa(shotNum)
		arg := callback.RequiredPlayers.SetArgs(shotStr).ToString()
		txt := hdlr.trans.Get(fmt.Sprintf("%s shot revolver", shotStr), translator.Config{})
		if settings.RequiredPlayers == shotNum {
			txt = fmt.Sprintf("🔫 %s", txt)
		}
		keyboard = append(keyboard, map[string]string{arg: txt})
	}

	hdlr.bot.SendInlineKeyboard(msg.Chat.ID, message, keyboard)
}
