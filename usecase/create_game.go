package usecase

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"gorm.io/gorm"
	"time"
)

var (
	ErrGameCooldown = errors.New("game cooldown")
)

type CreateGameUseCase struct {
	locker   service.Locker
	gameRepo domain.GameRepository
	userRepo domain.UserRepository
}

func NewCreateGameUseCase(
	locker service.Locker,
	gameRepo domain.GameRepository,
	userRepo domain.UserRepository,
) *CreateGameUseCase {
	return &CreateGameUseCase{
		gameRepo: gameRepo,
		userRepo: userRepo,
		locker:   locker,
	}
}

func (uc *CreateGameUseCase) Execute(chatID int64, ownerID int64) (*domain.Game, error) {
	locker := uc.locker.LockFor(fmt.Sprintf("game-start-%d", chatID))
	if !locker.TryLock() {
		return nil, service.ErrLockFailed
	}
	defer locker.Unlock()

	if uc.gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		game, err := uc.gameRepo.GetActiveForChat(chatID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGameCooldown
		}

		return game, nil
	}

	owner, err := uc.userRepo.Get(ownerID)
	if err != nil {
		return nil, err
	}

	game := &domain.Game{
		ID:        uuid.New(),
		ChatID:    chatID,
		OwnerID:   ownerID,
		Owner:     owner,
		CreatedAt: time.Now(),
	}

	if err := uc.gameRepo.Store(game); err != nil {
		return nil, err
	}

	return game, nil
}
