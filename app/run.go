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
	a.run(ctx, updates, maxConcurrentHandlers, shutdownTimeout)
}

func (a *App) run(ctx context.Context, updates <-chan tgbotapi.Update, maxConcurrent int, timeout time.Duration) {
	handlerCtx, cancelHandlers := context.WithCancel(context.WithoutCancel(ctx))
	defer cancelHandlers()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case sem <- struct{}{}:
		}

		select {
		case <-ctx.Done():
			<-sem
			break loop
		case update, ok := <-updates:
			if !ok {
				<-sem
				break loop
			}

			wg.Add(1)
			go func(update tgbotapi.Update) {
				defer wg.Done()
				defer func() { <-sem }()
				safeExecute(func() {
					a.HandleUpdate(handlerCtx, update)
				})
			}(update)
		}
	}

	waitForHandlers(&wg, timeout)
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
