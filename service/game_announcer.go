package service

import (
	"time"

	"github.com/taranovegor/naganbot/translator"
)

type announcerBot interface {
	SendMessage(chatID int64, text string)
	Kick(chatID int64, userID int64) error
}

// GameAnnouncer tells the chat how a game ended and kicks the victims.
type GameAnnouncer struct {
	bot          announcerBot
	trans        *translator.Translator
	messageDelay time.Duration
}

func NewGameAnnouncer(
	bot announcerBot,
	trans *translator.Translator,
	messageDelay time.Duration,
) *GameAnnouncer {
	return &GameAnnouncer{
		bot:          bot,
		trans:        trans,
		messageDelay: messageDelay,
	}
}

func (a *GameAnnouncer) AnnounceResult(chatID int64, report *HitReport) {
	for _, message := range a.trans.GetMany("play the game", translator.Config{}) {
		a.bot.SendMessage(chatID, message)
		time.Sleep(a.messageDelay)
	}

	isAtomic := report.BulletType == BulletAtomicType

	if isAtomic {
		a.bot.SendMessage(chatID, a.trans.Get("killed by atomic bullet", translator.Config{}))
	}

	for _, victim := range report.Victims {
		if !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
		}

		err := a.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("player is not kicked", translator.Config{}))
		}
	}
}
