package command

import (
	"context"
	"errors"
	"log"
	"math"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type resultAnnouncer interface {
	AnnounceResult(chatID int64, report *service.HitReport)
}

type joinHandler struct {
	bot            messenger
	createGameUC   *usecase.CreateGameUseCase
	joinGameUC     *usecase.JoinGameUseCase
	playGameUC     *usecase.PlayGameUseCase
	gunslingerRepo domain.GunslingerRepository
	announcer      resultAnnouncer
	trans          *translator.Translator
}

var _ Handler = (*joinHandler)(nil)

func NewJoinHandler(
	bot messenger,
	createGameUC *usecase.CreateGameUseCase,
	joinGameUC *usecase.JoinGameUseCase,
	playGameUC *usecase.PlayGameUseCase,
	gunslingerRepo domain.GunslingerRepository,
	announcer resultAnnouncer,
	trans *translator.Translator,
) Handler {
	return &joinHandler{
		bot:            bot,
		createGameUC:   createGameUC,
		joinGameUC:     joinGameUC,
		playGameUC:     playGameUC,
		gunslingerRepo: gunslingerRepo,
		announcer:      announcer,
		trans:          trans,
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
		if !errors.Is(err, usecase.ErrNotEnoughPlayers) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
			return
		}

		if game.Owner.ID != userID {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("joining the game", translator.Config{}))
		}
		if game.Mode == domain.GameModeCountdown {
			hdlr.sendCountdownStatus(chatID, game)
		}
		return
	}

	hdlr.announcer.AnnounceResult(chatID, hitReport)
}

func (hdlr *joinHandler) sendCountdownStatus(chatID int64, game *domain.Game) {
	gunslingers, err := hdlr.gunslingerRepo.GetByGameID(game.ID)
	if err != nil {
		log.Printf("failed to load gunslingers of game %s: %v", game.ID, err)
		return
	}

	startsAt, ok := game.StartsAt(gunslingers)
	if !ok {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("countdown waiting for players", translator.Config{
			Count: domain.CountdownMinPlayers - len(gunslingers),
		}))
		return
	}

	minutes := max(int(math.Ceil(time.Until(startsAt).Minutes())), 1)
	hdlr.bot.SendMessage(chatID, hdlr.trans.Get("countdown starts in", translator.Config{
		Count: minutes,
		Args:  map[string]string{"%max": strconv.Itoa(domain.CountdownMaxPlayers)},
	}))
}
