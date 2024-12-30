package command

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/taranovegor/naganbot/service/mock"
	"testing"
)

func TestForceHandler_Name(t *testing.T) {
	mockBot := new(mock.Bot)
	handler := NewForceHandler(mockBot)

	assert.Equal(t, "force", handler.Name())
}

func TestForceHandler_Execute(t *testing.T) {
	mockBot := new(mock.Bot)
	handler := NewForceHandler(mockBot)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 1111,
		},
		From: &tgbotapi.User{
			ID: 2222,
		},
		MessageID: 3333,
	}

	mockBot.On("Kick", msg.Chat.ID, msg.From.ID).Once().Return(nil)
	mockBot.On("DeleteMessage", msg.Chat.ID, msg.MessageID).Once().Return(nil)

	handler.Execute(msg)

	mockBot.AssertExpectations(t)
}
