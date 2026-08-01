package app

import (
	"context"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const shutdownTimeout = 10 * time.Second

func (a *App) Run(ctx context.Context, updates <-chan tgbotapi.Update) {
	handlerCtx, cancelHandlers := context.WithCancel(context.WithoutCancel(ctx))
	defer cancelHandlers()

	var wg sync.WaitGroup

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case update, ok := <-updates:
			if !ok {
				break loop
			}

			wg.Add(1)
			go func(update tgbotapi.Update) {
				defer wg.Done()
				safeExecute(func() {
					a.HandleUpdate(handlerCtx, update)
				})
			}(update)
		}
	}

	graceTimer := time.AfterFunc(shutdownTimeout, cancelHandlers)
	defer graceTimer.Stop()

	waitForHandlers(&wg, shutdownTimeout)
}

func safeExecute(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()
	fn()
}

func waitForHandlers(wg *sync.WaitGroup, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		log.Printf("timed out after %s waiting for in-flight update handlers", timeout)
	}
}
