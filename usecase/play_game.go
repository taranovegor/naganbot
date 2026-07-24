package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

var (
	ErrGameAlreadyPlayed = errors.New("game already played")
	ErrNotEnoughPlayers  = errors.New("not enough players")
)

type PlayGameUseCase struct {
	locker         service.Locker
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
	uow            domain.GameplayUnitOfWork
	nagan          *service.Nagan
}

func NewPlayGameUseCase(
	locker service.Locker,
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
	uow domain.GameplayUnitOfWork,
	nagan *service.Nagan,
) *PlayGameUseCase {
	return &PlayGameUseCase{
		locker:         locker,
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		uow:            uow,
		nagan:          nagan,
	}
}

func gameLockKey(gameID uuid.UUID) string {
	return fmt.Sprintf("play-game-%s", gameID.String())
}

func (uc *PlayGameUseCase) Execute(ctx context.Context, gameID uuid.UUID) (*service.HitReport, error) {
	locker := uc.locker.LockFor(gameLockKey(gameID))
	if !locker.TryLock() {
		return nil, service.ErrLockFailed
	}
	defer locker.Unlock()

	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	if game.IsPlayed() {
		return nil, ErrGameAlreadyPlayed
	}

	gunslingers, err := uc.gunslingerRepo.GetByGameID(game.ID)
	if err != nil {
		return nil, err
	}

	if len(gunslingers) < game.PlayersCount {
		return nil, ErrNotEnoughPlayers
	}

	report, err := uc.nagan.Shoot(ctx, game.ID, gunslingers)
	if err != nil {
		return nil, err
	}
	for _, victim := range report.Victims {
		victim.MarkAsShotHimself()
	}

	game.MarkAsPlayed(report.BulletType, report.ProofURL)
	if err := uc.uow.CommitPlayedGame(game, report.Victims); err != nil {
		return nil, err
	}

	return report, nil
}
