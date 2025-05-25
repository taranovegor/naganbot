package usecase

import (
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
	nagan          *service.Nagan
}

func NewPlayGameUseCase(
	locker service.Locker,
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
	nagan *service.Nagan,
) *PlayGameUseCase {
	return &PlayGameUseCase{
		locker:         locker,
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		nagan:          nagan,
	}
}

func (uc *PlayGameUseCase) Execute(gameID uuid.UUID) (*service.HitReport, error) {
	locker := uc.locker.LockFor(fmt.Sprintf("play-game-%d", gameID.ID()))
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

	report := uc.nagan.Shoot(gunslingers)
	for _, victim := range report.Victims {
		victim.MarkAsShotHimself()
	}

	game.MarkAsPlayed(report.BulletType)
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, err
	}

	if err := uc.gunslingerRepo.Update(report.Victims); err != nil {
		return nil, err
	}

	return report, nil
}
