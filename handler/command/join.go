package command

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"gorm.io/gorm/logger"
	"time"
)

type JoinHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	user       domain.UserRepository
	game       domain.GameRepository
	gunslinger domain.GunslingerRepository
	nagan      *service.Nagan
}

func NewJoinHandler(
	bot *service.Bot,
	trans *translator.Translator,
	userRepository domain.UserRepository,
	gameRepository domain.GameRepository,
	gunslingerRepository domain.GunslingerRepository,
	nagan *service.Nagan,
) Handler {
	return &JoinHandler{
		bot:        bot,
		trans:      trans,
		user:       userRepository,
		game:       gameRepository,
		gunslinger: gunslingerRepository,
		nagan:      nagan,
	}
}

func (hdlr JoinHandler) Name() string {
	return "join"
}

func (hdlr JoinHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	var message string

	activeGame, err := hdlr.game.GetActiveForChat(chatID)
	if errors.Is(err, logger.ErrRecordNotFound) {
		if hdlr.game.HasActiveOrCreatedTodayInChat(chatID) {
			message = hdlr.trans.Get("wait for game timeout", translator.Config{})
		} else {
			hdlr.game.Store(domain.NewGame(chatID, userID))

			message = hdlr.trans.Get("game creation", translator.Config{})
		}
	} else {
		if activeGame.OwnerID == userID || hdlr.gunslinger.IsPlayerExistsInGame(userID, activeGame.ID) {
			message = hdlr.trans.Get("player already in game", translator.Config{})
		} else {
			gunslinger := domain.NewGunslinger(activeGame.ID, userID)
			gunslinger.Player, _ = hdlr.user.Get(userID)
			hdlr.gunslinger.Store(gunslinger)
			activeGame.Gunslingers = append(activeGame.Gunslingers, gunslinger)
			gunslinger.Game = &activeGame

			message = hdlr.trans.Get("joining the game", translator.Config{})
		}
	}

	if activeGame.IsPlayed() || len(activeGame.Gunslingers) < 2 {
		hdlr.bot.SendMessage(chatID, message)

		return
	}

	for _, text := range hdlr.trans.GetMany("play the game", translator.Config{}) {
		message = text
		hdlr.bot.SendMessage(chatID, message)
		time.Sleep(time.Second)
	}

	report := hdlr.nagan.Shoot(activeGame.Gunslingers)
	for _, gunslinger := range report.Gunslingers {
		gunslinger.MarkAsShotHimself()
		hdlr.gunslinger.Update(gunslinger)
	}

	activeGame.MarkAsPlayed(report.BulletType)
	hdlr.game.Update(&activeGame)

	if service.BulletAtomicType == report.BulletType {
		message = hdlr.trans.Get("killed by atomic bullet", translator.Config{})
		hdlr.bot.SendMessage(chatID, message)
	}

	for _, gunslinger := range report.Gunslingers {
		message = hdlr.trans.Get("gunslinger killed", translator.Config{
			Args: map[string]string{"%gunslinger": gunslinger.Player.Mention()},
		})
		hdlr.bot.SendMessage(chatID, message)

		err = hdlr.bot.Kick(chatID, gunslinger.PlayerID)
		if err != nil && service.BulletAtomicType != report.BulletType {
			message = hdlr.trans.Get("player is not kicked", translator.Config{})
			hdlr.bot.SendMessage(chatID, message)
		}
	}
}
