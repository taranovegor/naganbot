package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
)

var (
	ErrPlayerAlreadyInGame = errors.New("player already in game")
)

type JoinGameUseCase struct {
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
	userRepo       domain.UserRepository
}

func NewJoinGameUseCase(
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
	userRepo domain.UserRepository,
) *JoinGameUseCase {
	return &JoinGameUseCase{
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		userRepo:       userRepo,
	}
}

func (uc *JoinGameUseCase) Execute(gameID uuid.UUID, userID int64) (*domain.Gunslinger, error) {
	if uc.gunslingerRepo.IsPlayerExistsInGame(userID, gameID) {
		return nil, ErrPlayerAlreadyInGame
	}

	user, err := uc.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	gunslinger := domain.NewGunslinger(gameID, user)

	if err := uc.gunslingerRepo.Store(gunslinger); err != nil {
		return nil, err
	}

	return gunslinger, nil
}
