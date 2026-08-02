package callback

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

var revolverOptions = []int{4, 5, 6, 7, 8, 9}

func RevolverKeyboard(selected int, trans *translator.Translator) service.Keyboard {
	var keyboard service.Keyboard
	for _, n := range revolverOptions {
		s := strconv.Itoa(n)
		arg := RequiredPlayers.SetArgs(s).ToString()
		txt := trans.Get(fmt.Sprintf("%s shot revolver", s), translator.Config{})
		if n == selected {
			txt = fmt.Sprintf("🔫 %s", txt)
		}
		keyboard = append(keyboard, []service.Button{{Data: arg, Text: txt}})
	}
	return keyboard
}

type requiredPlayers struct {
	chatRepo domain.ChatRepository
	bot      *service.Bot
	trans    *translator.Translator
}

var _ Handler = (*requiredPlayers)(nil)

func NewRequiredPlayers(
	chatRepo domain.ChatRepository,
	bot *service.Bot,
	trans *translator.Translator,
) Handler {
	return &requiredPlayers{
		chatRepo: chatRepo,
		bot:      bot,
		trans:    trans,
	}
}

func (hdlr *requiredPlayers) Pattern() Pattern {
	return RequiredPlayers
}

func (hdlr *requiredPlayers) Execute(_ context.Context, query *tgbotapi.CallbackQuery) {
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

	players, err := strconv.Atoi(RequiredPlayers.GetArg(query.Data, 1))
	if err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	chat, err := hdlr.chatRepo.Get(chatID)
	if err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	chat.Settings.RequiredPlayers = players
	if err := hdlr.chatRepo.UpdateSettings(&chat); err != nil {
		hdlr.bot.AnswerCallback(query.ID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	notification := fmt.Sprintf(
		"%s\n%s",
		hdlr.trans.Get("revolver has been replaced", translator.Config{Count: players}),
		hdlr.trans.Get("settings will be applied for next games", translator.Config{}),
	)
	hdlr.bot.AnswerCallback(query.ID, notification)

	keyboard := RevolverKeyboard(players, hdlr.trans)
	hdlr.bot.EditMessageReplyMarkup(chatID, query.Message.MessageID, keyboard)
}
