package app

import (
	"context"
	"errors"
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

func TestRunKeepsAnInFlightHandlerContextAliveAfterTheSignalContextIsCancelled(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	started := make(chan struct{})
	release := make(chan struct{})
	var errWhileInFlight error

	handler := &fakeCommandHandler{}
	handler.name = "join"
	handler.onExecute = func() {
		close(started)
		<-release
		errWhileInFlight = handler.capturedCtx.Err()
	}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      &fakeChatRepo{},
		Users:      &fakeUserRepo{},
		Commands:   command.NewRegistry("", handler),
		Callbacks:  callback.NewRegistry(),
	}

	updates := make(chan tgbotapi.Update)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		a.Run(ctx, updates)
		close(done)
	}()

	updates <- tgbotapi.Update{Message: commandMessage("/join")}
	<-started

	cancel()
	close(release)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected Run to return once the in-flight handler finished")
	}

	if handler.capturedCtx == nil {
		t.Fatal("expected the handler to capture a context")
	}
	if errors.Is(errWhileInFlight, context.Canceled) {
		t.Fatal("expected the in-flight handler's context to survive cancellation of the signal context, so it can finish within the grace period")
	}
}

func TestRunStopsDispatchingNewHandlersOnceTheSignalContextIsCancelled(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	release := make(chan struct{})
	first := &fakeCommandHandler{name: "join", onExecute: func() {
		<-release
	}}
	second := &fakeCommandHandler{name: "settings"}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      &fakeChatRepo{},
		Users:      &fakeUserRepo{},
		Commands:   command.NewRegistry("", first, second),
		Callbacks:  callback.NewRegistry(),
	}

	updates := make(chan tgbotapi.Update)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		a.Run(ctx, updates)
		close(done)
	}()

	updates <- tgbotapi.Update{Message: commandMessage("/join")}
	cancel()

	select {
	case updates <- tgbotapi.Update{Message: commandMessage("/settings")}:
		t.Fatal("expected the run loop to stop consuming updates once the context was cancelled")
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	<-done

	if second.executed {
		t.Fatal("expected no new handler to be dispatched after the context was cancelled")
	}
}
