package command

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
	"time"
)

type JoinHandler struct {
	Handler
	bot          *service.Bot
	createGameUC *usecase.CreateGameUseCase
	joinGameUC   *usecase.JoinGameUseCase
	playGameUC   *usecase.PlayGameUseCase
	trans        *translator.Translator
}

func NewJoinHandler(
	bot *service.Bot,
	createGameUC *usecase.CreateGameUseCase,
	joinGameUC *usecase.JoinGameUseCase,
	playGameUC *usecase.PlayGameUseCase,
	trans *translator.Translator,
) Handler {
	return &JoinHandler{
		bot:          bot,
		createGameUC: createGameUC,
		joinGameUC:   joinGameUC,
		playGameUC:   playGameUC,
		trans:        trans,
	}
}

func (h *JoinHandler) Name() string {
	return "join"
}

func (h *JoinHandler) Execute(msg *tgbotapi.Message) {
	chatID, userID := msg.Chat.ID, msg.From.ID
	game, err := h.createGameUC.Execute(chatID, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyInGame) {
			h.bot.SendMessage(chatID, h.trans.Get("player already in game", translator.Config{}))
		} else if errors.Is(err, usecase.ErrAlreadyPlayedToday) {
			h.bot.SendMessage(chatID, h.trans.Get("player already played today", translator.Config{}))
		}
		return
	}

	_, err = h.joinGameUC.Execute(game.ID, msg.From.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyInGame) {
			h.bot.SendMessage(chatID, h.trans.Get("player already in game", translator.Config{}))
		} else if errors.Is(err, usecase.ErrAlreadyPlayedToday) {
			h.bot.SendMessage(chatID, h.trans.Get("player already played today", translator.Config{}))
		}
		return
	}

	if game.Owner.ID == userID {
		h.bot.SendMessage(chatID, h.trans.Get("game creation", translator.Config{}))
	}

	hitReport, err := h.playGameUC.Execute(game.ID)
	if err != nil {
		if game.Owner.ID != userID && errors.Is(err, usecase.ErrNotEnoughPlayers) {
			h.bot.SendMessage(chatID, h.trans.Get("joining the game", translator.Config{}))
		}
		return
	}

	for _, message := range h.trans.GetMany("play the game", translator.Config{}) {
		h.bot.SendMessage(chatID, message)
		time.Sleep(time.Second)
	}

	isAtomic := hitReport.BulletType == service.BulletAtomicType

	if isAtomic {
		h.bot.SendMessage(chatID, h.trans.Get("killed by atomic bullet", translator.Config{}))
	}

	for _, victim := range hitReport.Victims {
		if !isAtomic {
			h.bot.SendMessage(chatID, h.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
		}

		err = h.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			h.bot.SendMessage(chatID, h.trans.Get("player is not kicked", translator.Config{}))
		}
	}
}
