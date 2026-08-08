package command

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/service"
)

type Handler interface {
	Name() string
	Execute(context.Context, *tgbotapi.Message)
}

type messenger interface {
	SendMessage(chatID int64, text string)
	SendInlineKeyboard(chatID int64, text string, keyboard service.Keyboard)
	Kick(chatID int64, userID int64) error
	DeleteMessage(chatID int64, messageID int)
}

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry(
	namePrefix string,
	handlers ...Handler,
) *Registry {
	reg := &Registry{
		handlers: make(map[string]Handler),
	}

	for _, hdlr := range handlers {
		name := fmt.Sprintf("%s%s", namePrefix, hdlr.Name())
		reg.handlers[name] = hdlr
	}

	return reg
}

func (reg Registry) Find(name string) (Handler, error) {
	hdlr, exists := reg.handlers[name]
	if !exists {
		return nil, fmt.Errorf("command %s not found", name)
	}

	return hdlr, nil
}
