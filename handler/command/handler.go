package command

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Name() string
	Execute(*tgbotapi.Message)
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
		return nil, errors.New(fmt.Sprintf("command %s not found", name))
	}

	return hdlr, nil
}
