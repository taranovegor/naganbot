package callback

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	RequiredPlayers Pattern = "required-players"
	PaginateJoined  Pattern = "paginate-joined"
)

type Pattern string

func (h Pattern) Name() string {
	return strings.Split(h.ToString(), "_")[0]
}

func (h Pattern) ToString() string {
	return string(h)
}

func (h Pattern) SetArgs(args ...string) Pattern {
	return Pattern(fmt.Sprintf("%s_%s", h.ToString(), strings.Join(args, "_")))
}

func (h Pattern) GetArg(withArgs string, argNum int) string {
	separated := strings.Split(strings.Replace(withArgs, h.ToString(), "", 1), "_")
	if 0 == argNum || len(separated) <= argNum {
		return ""
	}
	return separated[argNum]
}

type Handler interface {
	Pattern() Pattern
	Execute(*tgbotapi.CallbackQuery)
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
		return nil, errors.New(fmt.Sprintf("handler %s for query %s not found", query.Name(), query))
	}

	return hdlr, nil
}
