package callback

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
)

var revolverOptions = []int{4, 5, 6, 7, 8, 9}

func RevolverKeyboard(selected int, trans *translator.Translator) []map[string]string {
	var keyboard []map[string]string
	for _, n := range revolverOptions {
		s := strconv.Itoa(n)
		arg := RequiredPlayers.SetArgs(s).ToString()
		txt := trans.Get(fmt.Sprintf("%s shot revolver", s), translator.Config{})
		if n == selected {
			txt = fmt.Sprintf("🔫 %s", txt)
		}
		keyboard = append(keyboard, map[string]string{arg: txt})
	}
	return keyboard
}

type requiredPlayers struct {
	chatRepo domain.ChatRepository
	bot      *service.Bot
	trans    *translator.Translator
}

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

func (h *requiredPlayers) Pattern() Pattern {
	return RequiredPlayers
}

func (h *requiredPlayers) Execute(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	isAdmin, err := h.bot.IsAdmin(chatID, query.From.ID)
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if !isAdmin {
		notification := h.trans.Get("settings can be changed only by admins", translator.Config{})
		h.bot.AnswerCallback(query.ID, notification)
		return
	}

	players, err := strconv.Atoi(RequiredPlayers.GetArg(query.Data, 1))
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	chat, err := h.chatRepo.Get(chatID)
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	chat.Settings.RequiredPlayers = players
	h.chatRepo.Update(&chat)

	notification := fmt.Sprintf(
		"%s\n%s",
		h.trans.Get("revolver has been replaced", translator.Config{Count: players}),
		h.trans.Get("settings will be applied for next games", translator.Config{}),
	)
	h.bot.AnswerCallback(query.ID, notification)

	keyboard := RevolverKeyboard(players, h.trans)
	h.bot.EditMessageReplyMarkup(chatID, query.Message.MessageID, keyboard)
}
