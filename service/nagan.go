package service

import (
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/translator"
	"math/rand"
	"time"
)

type Nagan interface {
	Shot(gunslingers []*domain.Gunslinger) string
}

type NaganFactoryModel struct {
	Rarity int
	Nagan
}

type NaganFactory struct {
	models []NaganFactoryModel
}

func NewNaganFactory(models []NaganFactoryModel) *NaganFactory {
	return &NaganFactory{
		models: models,
	}
}

func (f NaganFactory) Create() Nagan {
	var rnd, acc int
	rnd = rand.Intn(100)

	for _, model := range f.models {
		if rnd > 100-model.Rarity-acc {
			return model.Nagan
		}

		acc += model.Rarity
	}

	return nil
}

type RegularNagan struct {
	Nagan
	trans *translator.Translator
}

func NewRegularNagan(trans *translator.Translator) Nagan {
	return &RegularNagan{
		trans: trans,
	}
}

func (ng RegularNagan) Shot(gunslingers []*domain.Gunslinger) string {
	gunslinger := gunslingers[rand.Intn(len(gunslingers))]
	gunslinger.MarkAsShotHimself()

	return ng.trans.Get("gunslinger killed", translator.Config{
		Args: map[string]string{"%gunslinger": gunslinger.Player.Mention()},
	})
}

type AtomicNagan struct {
	Nagan
	trans *translator.Translator
}

func NewAtomicNagan(trans *translator.Translator) Nagan {
	return &AtomicNagan{
		trans: trans,
	}
}

func (ng AtomicNagan) Shot(gunslingers []*domain.Gunslinger) string {
	for _, g := range gunslingers {
		g.MarkAsShotHimself()
	}

	return ng.trans.Get("killed by atomic bullet", translator.Config{})
}

type Croupier struct {
	nagan                *NaganFactory
	bot                  *Bot
	trans                *translator.Translator
	gameRepository       domain.GameRepository
	gunslingerRepository domain.GunslingerRepository
}

func NewCroupier(
	nagan *NaganFactory,
	bot *Bot,
	trans *translator.Translator,
	gameRepository domain.GameRepository,
	gunslingerRepository domain.GunslingerRepository,
) *Croupier {
	return &Croupier{
		nagan:                nagan,
		bot:                  bot,
		trans:                trans,
		gameRepository:       gameRepository,
		gunslingerRepository: gunslingerRepository,
	}
}

func (c Croupier) Play(game *domain.Game) {
	chatID := game.ChatID

	for _, text := range c.trans.GetMany("play the game", translator.Config{}) {
		c.bot.SendMessage(chatID, text)
		time.Sleep(time.Second)
	}

	msg := c.nagan.Create().Shot(game.Gunslingers)
	c.bot.SendMessage(chatID, msg)
	for _, gunslinger := range game.Gunslingers {
		if gunslinger.ShotHimself {
			c.gunslingerRepository.Update(gunslinger)
			c.bot.Kick(chatID, gunslinger.PlayerID)
		}
	}

	game.MarkAsPlayed()
	c.gameRepository.Update(game)
}
