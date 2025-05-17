package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
)

type StatHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	gunslinger domain.GunslingerRepository
}

func NewStatHandler(
	bot *service.Bot,
	trans *translator.Translator,
	gunslinger domain.GunslingerRepository,
) Handler {
	return &StatHandler{
		bot:        bot,
		trans:      trans,
		gunslinger: gunslinger,
	}
}

func (hdlr StatHandler) Name() string {
	return "stat"
}

func (hdlr StatHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	numberOfGames := hdlr.gunslinger.CountNumberOfPlayerGamesInChat(userID, chatID)
	numberOfShotHimself := hdlr.gunslinger.CountNumberOfSelfShotsInChat(userID, chatID)

	hdlr.bot.SendMessage(chatID, hdlr.trans.Get("user game statistics", translator.Config{
		Args: map[string]string{
			"%games": strconv.FormatInt(numberOfGames, 10),
			"%shots": strconv.FormatInt(numberOfShotHimself, 10),
		},
		Count: int(numberOfShotHimself),
	}))
}
