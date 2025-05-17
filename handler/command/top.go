package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
)

type TopHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	user       domain.UserRepository
	gunslinger domain.GunslingerRepository
}

func NewTopHandler(
	bot *service.Bot,
	trans *translator.Translator,
	user domain.UserRepository,
	gunslinger domain.GunslingerRepository,
) Handler {
	return &TopHandler{
		bot:        bot,
		trans:      trans,
		user:       user,
		gunslinger: gunslinger,
	}
}

func (hdlr TopHandler) Name() string {
	return "top"
}

func (hdlr TopHandler) Execute(msg *tgbotapi.Message) {
	var players []domain.GunslingerTopShotPlayer
	chatID := msg.Chat.ID
	if year, err := strconv.Atoi(msg.CommandArguments()); err == nil {
		players, err = hdlr.gunslinger.GetTopShopPlayersByYearInChat(chatID, year)
	} else {
		players, err = hdlr.gunslinger.GetTopShotPlayersInChat(chatID)
	}

	var message string
	if len(players) == 0 {
		message = hdlr.trans.Get("top is not determined", translator.Config{})
	} else {
		message = hdlr.trans.Get("top players by games", translator.Config{
			Args: map[string]string{"%number": strconv.Itoa(len(players))},
		})
		for i, player := range players {
			// todo: need to reduce the number of queries
			user, err := hdlr.user.Get(player.PlayerId)
			if err != nil {
				continue
			}

			message += "\n" + hdlr.trans.Get("top game player", translator.Config{
				Args: map[string]string{
					"%i":     strconv.Itoa(i + 1),
					"%user":  user.Name(),
					"%times": strconv.Itoa(player.Times),
				},
				Count: player.Times,
			})
		}
	}

	hdlr.bot.SendMessage(chatID, message)
}
