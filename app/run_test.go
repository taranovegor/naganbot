package app

import (
	"context"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/translator"
)

func TestRunWaitsForInFlightHandlersAfterTheUpdatesChannelCloses(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	handler := &fakeCommandHandler{name: "join"}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      &fakeChatRepo{},
		Users:      &fakeUserRepo{},
		Commands:   command.NewRegistry("", handler),
		Callbacks:  callback.NewRegistry(),
	}

	updates := make(chan tgbotapi.Update, 1)
	updates <- tgbotapi.Update{Message: commandMessage("/join")}
	close(updates)

	a.Run(context.Background(), updates)

	if !handler.executed {
		t.Fatal("expected the handler launched before the channel closed to have finished by the time Run returns")
	}
}

func TestRunReturnsImmediatelyWhenThereAreNoUpdates(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      &fakeChatRepo{},
		Users:      &fakeUserRepo{},
		Commands:   command.NewRegistry(""),
		Callbacks:  callback.NewRegistry(),
	}

	updates := make(chan tgbotapi.Update)
	close(updates)

	done := make(chan struct{})
	go func() {
		a.Run(context.Background(), updates)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected Run to return quickly when there is nothing to wait for")
	}
}
