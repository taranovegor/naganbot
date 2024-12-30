package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/service"
)

type Bot struct {
	mock.Mock
	service.Bot
}

func (m *Bot) SendMessage(chatID int64, text string) error {
	args := m.Called(chatID, text)
	return args.Error(0)
}

func (m *Bot) DeleteMessage(chatID int64, messageID int) error {
	args := m.Called(chatID, messageID)
	return args.Error(0)
}

func (m *Bot) Ban(chatID int64, userID int64, untilDate int64) error {
	args := m.Called(chatID, userID, untilDate)
	return args.Error(0)
}

func (m *Bot) Kick(chatID int64, userID int64) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}
