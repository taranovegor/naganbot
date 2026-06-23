package command

import (
	"context"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type DuelHandler struct {
	bot      messenger
	trans    *translator.Translator
	userRepo domain.UserRepository
	duelRepo domain.DuelRepository
}

func NewDuelHandler(
	bot messenger,
	trans *translator.Translator,
	userRepo domain.UserRepository,
	duelRepo domain.DuelRepository,
) Handler {
	return &DuelHandler{
		bot:      bot,
		trans:    trans,
		userRepo: userRepo,
		duelRepo: duelRepo,
	}
}

func (h *DuelHandler) Name() string {
	return "duel"
}

func (h *DuelHandler) Execute(_ context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	challengerID := msg.From.ID

	var opponent domain.User

	if msg.ReplyToMessage != nil {
		reply := msg.ReplyToMessage.From
		if reply.IsBot {
			h.bot.SendMessage(chatID, h.trans.Get("duel is bot", translator.Config{}))
			return
		}
		opponent = domain.User{
			ID:        reply.ID,
			FirstName: reply.FirstName,
			LastName:  reply.LastName,
			Username:  reply.UserName,
		}
	} else {
		args := msg.CommandArguments()
		username := strings.TrimPrefix(strings.TrimSpace(args), "@")
		if username == "" {
			h.bot.SendMessage(chatID, h.trans.Get("duel no target", translator.Config{}))
			return
		}

		user, err := h.userRepo.GetByUsername(username)
		if err != nil {
			h.bot.SendMessage(chatID, h.trans.Get("duel unknown user", translator.Config{}))
			return
		}
		opponent = user
	}

	if opponent.ID == challengerID {
		err := h.bot.Kick(chatID, challengerID)
		if err != nil {
			h.bot.SendMessage(chatID, h.trans.Get("duel self", translator.Config{}))
		}
		return
	}

	if !h.bot.IsChatMember(chatID, opponent.ID) {
		h.bot.SendMessage(chatID, h.trans.Get("duel not in chat", translator.Config{}))
		return
	}

	if h.bot.IsBot(chatID, opponent.ID) {
		h.bot.SendMessage(chatID, h.trans.Get("duel is bot", translator.Config{}))
		return
	}

	pending, err := h.duelRepo.HasPendingFrom(chatID, challengerID, time.Now().Add(-domain.DuelTTL))
	if err != nil {
		h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}
	if pending {
		h.bot.SendMessage(chatID, h.trans.Get("duel pending", translator.Config{}))
		return
	}

	challenger, err := h.userRepo.Get(challengerID)
	if err != nil {
		h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	duel := domain.NewDuel(chatID, challenger, opponent)
	if err := h.duelRepo.Store(duel); err != nil {
		h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	message := h.trans.Get("duel challenge", translator.Config{
		Args: map[string]string{
			"%challenger": challenger.Mention(),
			"%opponent":   opponent.Mention(),
		},
	})

	acceptArg := callback.DuelResponse.SetArgs(duel.ID.String(), "accept").ToString()
	declineArg := callback.DuelResponse.SetArgs(duel.ID.String(), "decline").ToString()

	keyboard := service.Keyboard{
		{
			{Data: acceptArg, Text: h.trans.Get("duel accept button", translator.Config{})},
			{Data: declineArg, Text: h.trans.Get("duel decline button", translator.Config{})},
		},
	}

	h.bot.SendInlineKeyboard(chatID, message, keyboard)
}
