package view

import (
	"fmt"
	"strconv"

	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type JoinedMessage struct {
	Text     string
	Keyboard service.InlineKeyboard
}

type JoinedMessageFactory struct {
	translator *translator.Translator
}

func NewJoinedMessageFactory(
	trans *translator.Translator,
) *JoinedMessageFactory {
	return &JoinedMessageFactory{
		translator: trans,
	}
}

func (f *JoinedMessageFactory) Create(game *domain.Game) JoinedMessage {
	if game == nil {
		return JoinedMessage{
			Text:     f.translator.Get("active game not found", translator.Config{}),
			Keyboard: nil,
		}
	}

	text := f.translator.Get("joined the game", translator.Config{
		Args: map[string]string{"%date": game.CreatedAt.Format(config.DateFormat)},
	})

	if game.IsPlayed() && game.BulletType == service.BulletAtomicType {
		text = "☢️ " + text
	}

	for i, gunslinger := range game.Gunslingers {
		line := f.translator.Get("game join list item", translator.Config{
			Args: map[string]string{
				"%num":        strconv.Itoa(i + 1),
				"%gunslinger": gunslinger.Player.Name(),
			},
		})

		if game.OwnerID == gunslinger.PlayerID {
			line += " " + f.translator.Get("owner of the game", translator.Config{})
		}

		if gunslinger.ShotHimself {
			if game.OwnerID == gunslinger.PlayerID {
				line += ","
			}
			line += " " + f.translator.Get("shot in game", translator.Config{})
		}

		text += "\n" + line
	}

	gameID := game.ID.String()
	keyboard := service.InlineKeyboard{
		{
			fmt.Sprintf("paginate-joined_prev_%s", gameID): "⬅️",
			fmt.Sprintf("paginate-joined_next_%s", gameID): "➡️",
		},
	}

	return JoinedMessage{
		Text:     text,
		Keyboard: keyboard,
	}
}
