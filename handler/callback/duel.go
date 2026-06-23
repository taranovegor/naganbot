package callback

import (
	"context"
	"errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type duelResponse struct {
	bot     messenger
	trans   *translator.Translator
	respond *usecase.RespondDuelUseCase
}

var _ Handler = (*duelResponse)(nil)

func NewDuelResponse(
	bot messenger,
	trans *translator.Translator,
	respond *usecase.RespondDuelUseCase,
) Handler {
	return &duelResponse{bot: bot, trans: trans, respond: respond}
}

func (h *duelResponse) Pattern() Pattern {
	return DuelResponse
}

func (h *duelResponse) Execute(ctx context.Context, query *tgbotapi.CallbackQuery) {
	duelID, err := uuid.Parse(DuelResponse.GetArg(query.Data, 1))
	if err != nil {
		h.fail(query)
		return
	}

	switch DuelResponse.GetArg(query.Data, 2) {
	case "accept":
		h.accept(ctx, query, duelID)
	case "decline":
		h.decline(query, duelID)
	}
}

func (h *duelResponse) accept(ctx context.Context, query *tgbotapi.CallbackQuery, duelID uuid.UUID) {
	outcome, err := h.respond.Accept(ctx, duelID, query.From.ID, time.Now())
	if err != nil {
		h.handleError(query, err)
		return
	}

	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	duel := outcome.Duel

	h.bot.AnswerCallback(query.ID, "")
	h.bot.EditMessageText(chatID, messageID, h.trans.Get("duel accepted", translator.Config{
		Args: map[string]string{
			"%challenger": duel.Challenger.Mention(),
			"%opponent":   duel.Opponent.Mention(),
		},
	}))

	result := h.trans.Get("duel result", translator.Config{
		Args: map[string]string{"%loser": outcome.Loser.Mention()},
	})
	result += "\n\n" + h.trans.Get("proof link", translator.Config{
		Args: map[string]string{"%url": outcome.ProofURL},
	})
	h.bot.ReplyMessage(chatID, messageID, result)

	if err := h.bot.Kick(chatID, outcome.Loser.ID); err != nil {
		log.Printf("failed to kick duel loser %d in chat %d: %v", outcome.Loser.ID, chatID, err)
		h.bot.ReplyMessage(chatID, messageID, h.trans.Get("duel kick failed", translator.Config{}))
	}
}

func (h *duelResponse) decline(query *tgbotapi.CallbackQuery, duelID uuid.UUID) {
	duel, err := h.respond.Decline(duelID, query.From.ID, time.Now())
	if err != nil {
		h.handleError(query, err)
		return
	}

	h.bot.AnswerCallback(query.ID, "")
	h.bot.EditMessageText(query.Message.Chat.ID, query.Message.MessageID, h.trans.Get("duel declined", translator.Config{
		Args: map[string]string{
			"%challenger": duel.Challenger.Mention(),
			"%opponent":   duel.Opponent.Mention(),
		},
	}))
}

func (h *duelResponse) handleError(query *tgbotapi.CallbackQuery, err error) {
	switch {
	case errors.Is(err, usecase.ErrDuelNotOpponent):
		h.bot.AnswerCallback(query.ID, h.trans.Get("duel not your", translator.Config{}))
	case errors.Is(err, usecase.ErrDuelClosed), errors.Is(err, service.ErrLockNotAcquired):
		h.bot.AnswerCallback(query.ID, "")
	case errors.Is(err, usecase.ErrDuelExpired):
		h.bot.AnswerCallback(query.ID, "")
		h.bot.EditMessageText(query.Message.Chat.ID, query.Message.MessageID, h.trans.Get("duel expired", translator.Config{}))
	default:
		log.Printf("failed to respond to duel: %v", err)
		h.fail(query)
	}
}

func (h *duelResponse) fail(query *tgbotapi.CallbackQuery) {
	h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
}
