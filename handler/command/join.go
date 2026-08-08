package command

import (
	"context"
	"errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type joinHandler struct {
	bot          messenger
	createGameUC *usecase.CreateGameUseCase
	joinGameUC   *usecase.JoinGameUseCase
	playGameUC   *usecase.PlayGameUseCase
	trans        *translator.Translator
}

var _ Handler = (*joinHandler)(nil)

func NewJoinHandler(
	bot messenger,
	createGameUC *usecase.CreateGameUseCase,
	joinGameUC *usecase.JoinGameUseCase,
	playGameUC *usecase.PlayGameUseCase,
	trans *translator.Translator,
) Handler {
	return &joinHandler{
		bot:          bot,
		createGameUC: createGameUC,
		joinGameUC:   joinGameUC,
		playGameUC:   playGameUC,
		trans:        trans,
	}
}

func (hdlr *joinHandler) Name() string {
	return "join"
}

func (hdlr *joinHandler) Execute(ctx context.Context, msg *tgbotapi.Message) {
	chatID, userID := msg.Chat.ID, msg.From.ID
	game, err := hdlr.createGameUC.Execute(chatID, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrGameCooldown) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("wait for game timeout", translator.Config{}))
		}
		return
	}

	_, err = hdlr.joinGameUC.Execute(game.ID, msg.From.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerAlreadyInGame) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("player already in game", translator.Config{}))
		}
		return
	}

	if game.Owner.ID == userID {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("game creation", translator.Config{}))
	}

	hitReport, err := hdlr.playGameUC.Execute(ctx, game.ID)
	if err != nil {
		log.Printf("failed to play game %s: %v", game.ID, err)
		if game.Owner.ID != userID && errors.Is(err, usecase.ErrNotEnoughPlayers) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("joining the game", translator.Config{}))
		} else if !errors.Is(err, usecase.ErrNotEnoughPlayers) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
		}
		return
	}

	for _, message := range hdlr.trans.GetMany("play the game", translator.Config{}) {
		hdlr.bot.SendMessage(chatID, message)
		time.Sleep(time.Second)
	}

	isAtomic := hitReport.BulletType == service.BulletAtomicType

	if isAtomic {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("killed by atomic bullet", translator.Config{}))
	}

	for _, victim := range hitReport.Victims {
		if !isAtomic {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
		}

		err = hdlr.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("player is not kicked", translator.Config{}))
		}
	}
}
