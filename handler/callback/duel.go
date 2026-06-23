package callback

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type duelResponse struct {
	bot      *service.Bot
	trans    *translator.Translator
	duelRepo domain.DuelRepository
	drand    *drand.Client
}

func NewDuelResponse(
	bot *service.Bot,
	trans *translator.Translator,
	duelRepo domain.DuelRepository,
	drand *drand.Client,
) Handler {
	return &duelResponse{
		bot:      bot,
		trans:    trans,
		duelRepo: duelRepo,
		drand:    drand,
	}
}

func (h *duelResponse) Pattern() Pattern {
	return DuelResponse
}

func (h *duelResponse) Execute(query *tgbotapi.CallbackQuery) {
	userID := query.From.ID

	duelID, err := uuid.Parse(DuelResponse.GetArg(query.Data, 1))
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	action := DuelResponse.GetArg(query.Data, 2)

	duel, err := h.duelRepo.GetByID(duelID)
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if duel.IsPlayed() {
		h.bot.AnswerCallback(query.ID, "")
		return
	}

	if userID != duel.OpponentID {
		h.bot.AnswerCallback(query.ID, h.trans.Get("duel not your", translator.Config{}))
		return
	}

	switch action {
	case "accept":
		h.accept(query, duel)
	case "decline":
		h.decline(query, duel)
	}
}

func (h *duelResponse) accept(query *tgbotapi.CallbackQuery, duel *domain.Duel) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	h.bot.AnswerCallback(query.ID, "")

	h.bot.EditMessageText(chatID, messageID, h.trans.Get("duel accepted", translator.Config{
		Args: map[string]string{
			"%challenger": duel.Challenger.Mention(),
			"%opponent":   duel.Opponent.Mention(),
		},
	}))

	loserID, proofURL, err := h.shoot(duel)
	if err != nil {
		h.bot.ReplyMessage(chatID, messageID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	var loser domain.User
	if loserID == duel.ChallengerID {
		loser = duel.Challenger
	} else {
		loser = duel.Opponent
	}

	duel.MarkAsPlayed(loserID, proofURL)
	h.duelRepo.Update(duel)

	result := h.trans.Get("duel result", translator.Config{
		Args: map[string]string{"%loser": loser.Mention()},
	})
	result += "\n\n" + h.trans.Get("proof link", translator.Config{
		Args: map[string]string{"%url": proofURL},
	})
	h.bot.ReplyMessage(chatID, messageID, result)

	h.bot.Kick(chatID, loserID)
}

func (h *duelResponse) decline(query *tgbotapi.CallbackQuery, duel *domain.Duel) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	h.bot.AnswerCallback(query.ID, "")

	h.bot.EditMessageText(chatID, messageID, h.trans.Get("duel declined", translator.Config{
		Args: map[string]string{
			"%challenger": duel.Challenger.Mention(),
			"%opponent":   duel.Opponent.Mention(),
		},
	}))
}

func (h *duelResponse) shoot(duel *domain.Duel) (loserID int64, proofURL string, err error) {
	beacon, err := h.drand.GetLatest(context.TODO())
	if err != nil {
		return 0, "", fmt.Errorf("failed to get drand beacon: %w", err)
	}

	randomnessBytes, err := hex.DecodeString(beacon.Randomness)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode randomness: %w", err)
	}

	hash := sha256.Sum256(append(randomnessBytes, []byte(duel.ID.String())...))
	result := binary.BigEndian.Uint64(hash[4:12]) % 2

	if result == 0 {
		loserID = duel.ChallengerID
	} else {
		loserID = duel.OpponentID
	}

	return loserID, h.drand.ProofURL(beacon.Round, duel.ID), nil
}
