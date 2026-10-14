package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

type StartedGame struct {
	ChatID int64
	Report *service.HitReport
}

// StartDueGamesUseCase plays every countdown game whose time has come. It keeps
// no state: what is due depends only on the stored join times, so a restart
// cannot lose or repeat a start.
type StartDueGamesUseCase struct {
	gameRepo domain.GameRepository
	playGame *PlayGameUseCase
}

func NewStartDueGamesUseCase(
	gameRepo domain.GameRepository,
	playGame *PlayGameUseCase,
) *StartDueGamesUseCase {
	return &StartDueGamesUseCase{
		gameRepo: gameRepo,
		playGame: playGame,
	}
}

// Execute returns the games it played and one joined error for the games that
// failed. If someone else played a game first, that is not an error.
func (uc *StartDueGamesUseCase) Execute(ctx context.Context, now time.Time) ([]StartedGame, error) {
	games, err := uc.gameRepo.GetActiveCountdownGames(now.Add(-domain.LobbyIdleTTL))
	if err != nil {
		return nil, err
	}

	var (
		started []StartedGame
		errs    []error
	)
	for _, game := range games {
		if !game.ShouldStart(game.Gunslingers, now) {
			continue
		}

		report, err := uc.playGame.Execute(ctx, game.ID)
		switch {
		case err == nil:
			started = append(started, StartedGame{ChatID: game.ChatID, Report: report})
		case isGameTakenCareOf(err):
		default:
			errs = append(errs, err)
		}
	}

	return started, errors.Join(errs...)
}

func isGameTakenCareOf(err error) bool {
	return errors.Is(err, service.ErrLockNotAcquired) ||
		errors.Is(err, ErrGameAlreadyPlayed) ||
		errors.Is(err, ErrNotEnoughPlayers)
}
