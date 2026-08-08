package command

import (
	"context"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type historyHandler struct {
	bot   messenger
	trans *translator.Translator
	game  domain.GameRepository
}

var _ Handler = (*historyHandler)(nil)

func NewLogHandler(
	bot messenger,
	trans *translator.Translator,
	game domain.GameRepository,
) Handler {
	return &historyHandler{
		bot:   bot,
		trans: trans,
		game:  game,
	}
}

func (hdlr *historyHandler) Name() string {
	return "history"
}

func (hdlr *historyHandler) Execute(_ context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	limit := 10
	if arg := msg.CommandArguments(); len(arg) > 0 {
		if n, err := strconv.Atoi(arg); err == nil && n > 0 {
			limit = n
			if n > 50 {
				limit = 50
			}
		}
	}

	games, err := hdlr.game.GetLatestGamesInChat(chatID, limit)
	if err != nil {
		log.Printf("failed to get game log: %v", err)
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if len(games) == 0 {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("active game not found", translator.Config{}))
		return
	}

	message := hdlr.trans.Get("game log header", translator.Config{
		Args: map[string]string{"%number": strconv.Itoa(len(games))},
	})
	for _, game := range games {
		username := findLoser(game)
		message += "\n" + hdlr.trans.Get("game log item", translator.Config{
			Args: map[string]string{
				"%date":     game.PlayedAt.Time.Format("02.01"),
				"%username": username,
			},
		})
	}

	hdlr.bot.SendMessage(chatID, message)
}

func findLoser(game domain.Game) string {
	if game.BulletType == service.BulletAtomicType {
		var names []string
		for _, gs := range game.Gunslingers {
			names = append(names, gs.Player.Name())
		}

		return "☢️ " + strings.Join(names, ", ")
	}

	for _, gs := range game.Gunslingers {
		if gs.ShotHimself {
			return gs.Player.Name()
		}
	}

	return game.Owner.Name()
}
