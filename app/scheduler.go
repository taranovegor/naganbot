package app

import (
	"context"
	"log"
	"time"

	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/usecase"
)

const schedulerInterval = time.Minute

type dueGamesStarter interface {
	Execute(ctx context.Context, now time.Time) ([]usecase.StartedGame, error)
}

type resultAnnouncer interface {
	AnnounceResult(chatID int64, report *service.HitReport)
}

// RunScheduler starts the countdown games that are due. It checks once at
// startup, to catch up after a restart, and then on every tick until ctx is done.
func (a *App) RunScheduler(ctx context.Context) {
	runScheduler(ctx, schedulerInterval, a.dueGames, a.announcer)
}

func runScheduler(ctx context.Context, interval time.Duration, starter dueGamesStarter, announcer resultAnnouncer) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		startDueGames(ctx, time.Now(), starter, announcer)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func startDueGames(ctx context.Context, now time.Time, starter dueGamesStarter, announcer resultAnnouncer) {
	started, err := starter.Execute(ctx, now)
	if err != nil {
		log.Printf("failed to start due games: %v", err)
	}

	for _, game := range started {
		announcer.AnnounceResult(game.ChatID, game.Report)
	}
}
