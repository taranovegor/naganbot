package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
)

type JoinedHandler struct {
	Handler
	bot   *service.Bot
	trans *translator.Translator
	game  domain.GameRepository
}

func NewJoinedHandler(
	bot *service.Bot,
	trans *translator.Translator,
	game domain.GameRepository,
) Handler {
	return &JoinedHandler{
		bot:   bot,
		trans: trans,
		game:  game,
	}
}

func (hdlr JoinedHandler) Name() string {
	return "joined"
}

func (hdlr JoinedHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	var message string

	game, err := hdlr.game.GetLatestForChat(chatID)
	if err != nil {
		message = hdlr.trans.Get("active game not found", translator.Config{})
	} else {
		message = hdlr.trans.Get("joined the game", translator.Config{
			Args: map[string]string{"%date": game.CreatedAt.Format(config.DateFormat)},
		})
		if game.IsPlayed() && game.BulletType == service.BulletAtomicType {
			message = "☢️ " + message
		}
		for i, gunslinger := range game.Gunslingers {
			message += "\n" + hdlr.trans.Get("game join list item", translator.Config{
				Args: map[string]string{"%num": strconv.Itoa(i + 1), "%gunslinger": gunslinger.Player.Name()},
			})

			if game.OwnerID == gunslinger.PlayerID {
				message += " " + hdlr.trans.Get("owner of the game", translator.Config{})
			}

			if gunslinger.ShotHimself {
				if game.OwnerID == gunslinger.PlayerID {
					message += ","
				}

				message += " " + hdlr.trans.Get("shot in game", translator.Config{})
			}
		}
	}

	hdlr.bot.SendMessage(chatID, message)
}
