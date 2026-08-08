package command

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestForceHandlerExecuteKicksTheSenderAndDeletesTheirMessage(t *testing.T) {
	bot := &joinTestBot{}
	hdlr := &forceHandler{bot: bot}

	msg := &tgbotapi.Message{
		Chat:      &tgbotapi.Chat{ID: 100},
		From:      &tgbotapi.User{ID: 7},
		MessageID: 55,
	}

	hdlr.Execute(context.Background(), msg)

	kicked := bot.kickedPlayers()
	if len(kicked) != 1 || kicked[0] != 7 {
		t.Fatalf("expected user 7 to be kicked, got %v", kicked)
	}
	if !bot.deleteCalled || bot.deletedMessage != 55 {
		t.Fatalf("expected message 55 to be deleted, got called=%v id=%d", bot.deleteCalled, bot.deletedMessage)
	}
}

func TestForceHandlerExecuteStillDeletesTheMessageWhenTheKickFails(t *testing.T) {
	bot := &joinTestBot{kickErr: errors.New("telegram: kick failed")}
	hdlr := &forceHandler{bot: bot}

	msg := &tgbotapi.Message{
		Chat:      &tgbotapi.Chat{ID: 100},
		From:      &tgbotapi.User{ID: 7},
		MessageID: 55,
	}

	hdlr.Execute(context.Background(), msg)

	if !bot.deleteCalled {
		t.Fatal("expected the message to be deleted even though the kick failed")
	}
}
