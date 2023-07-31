package command

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Name() string
	Execute(*tgbotapi.Message) error
}

type Registry struct {
	commands map[string]Command
}

func NewRegistry(
	namePrefix string,
	commands ...Command,
) *Registry {
	reg := &Registry{
		commands: make(map[string]Command),
	}

	for _, cmd := range commands {
		name := fmt.Sprintf("%s%s", namePrefix, cmd.Name())
		reg.commands[name] = cmd
	}

	return reg
}

func (reg Registry) Find(name string) (Command, error) {
	cmd, exists := reg.commands[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("command %s not found", name))
	}

	return cmd, nil
}
