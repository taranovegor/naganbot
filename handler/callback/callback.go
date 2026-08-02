package callback

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	RequiredPlayers Pattern = "required-players"
)

type Pattern string

func (p Pattern) Name() string {
	return strings.Split(p.ToString(), "_")[0]
}

func (p Pattern) ToString() string {
	return string(p)
}

func (p Pattern) SetArgs(args ...string) Pattern {
	return Pattern(fmt.Sprintf("%s_%s", p.ToString(), strings.Join(args, "_")))
}

func (p Pattern) GetArg(withArgs string, argNum int) string {
	separated := strings.Split(strings.Replace(withArgs, p.ToString(), "", 1), "_")
	if 0 == argNum || len(separated) <= argNum {
		return ""
	}
	return separated[argNum]
}

type Handler interface {
	Pattern() Pattern
	Execute(context.Context, *tgbotapi.CallbackQuery)
}

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry(
	handlers ...Handler,
) *Registry {
	reg := &Registry{
		handlers: make(map[string]Handler),
	}

	for _, hdlr := range handlers {
		name := hdlr.Pattern().Name()
		reg.handlers[name] = hdlr
	}

	return reg
}

func (reg Registry) Find(query Pattern) (Handler, error) {
	hdlr, exists := reg.handlers[query.Name()]
	if !exists {
		return nil, fmt.Errorf("handler %s for query %s not found", query.Name(), query)
	}

	return hdlr, nil
}
