package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
)

var (
	ErrUserAlreadyInGame  = errors.New("player already in game")
	ErrAlreadyPlayedToday = errors.New("player already played today")
)

type JoinGameUseCase struct {
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
}

func NewJoinGameUseCase(
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
) *JoinGameUseCase {
	return &JoinGameUseCase{
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
	}
}

func (uc *JoinGameUseCase) Execute(gameID uuid.UUID, userID int64) (*domain.Gunslinger, error) {
	if uc.gunslingerRepo.IsUserInGame(userID, gameID) {
		return nil, ErrUserAlreadyInGame
	}

	if uc.gunslingerRepo.HasPlayedToday(userID) {
		return nil, ErrAlreadyPlayedToday
	}

	gunslinger := domain.NewGunslinger(
		gameID,
		userID,
	)

	err := uc.gunslingerRepo.Store(gunslinger)
	if err != nil {
		return nil, err
	}

	return gunslinger, nil
}
