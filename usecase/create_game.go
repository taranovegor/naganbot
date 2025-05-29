package usecase

import (
	"errors"
	"fmt"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"gorm.io/gorm"
)

type CreateGameUseCase struct {
	locker         service.Locker
	chatRepo       domain.ChatRepository
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
}

func NewCreateGameUseCase(
	locker service.Locker,
	chatRepo domain.ChatRepository,
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
) *CreateGameUseCase {
	return &CreateGameUseCase{
		locker:         locker,
		chatRepo:       chatRepo,
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
	}
}

func (uc *CreateGameUseCase) Execute(chatID int64, ownerID int64) (*domain.Game, error) {
	locker := uc.locker.LockFor(fmt.Sprintf("game-start-%d", chatID))
	if !locker.TryLock() {
		return nil, service.ErrLockFailed
	}
	defer locker.Unlock()

	game, err := uc.gameRepo.GetActiveInChat(chatID)
	if err == nil {
		if uc.gunslingerRepo.IsUserInGame(ownerID, game.ID) {
			return nil, ErrUserAlreadyInGame
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if uc.gunslingerRepo.HasPlayedToday(ownerID) {
		return nil, ErrAlreadyPlayedToday
	}

	chat, err := uc.chatRepo.Get(chatID)
	if err != nil {
		return nil, err
	}

	game = domain.NewGame(
		chatID,
		ownerID,
		chat.Settings.RequiredPlayers,
	)

	if err := uc.gameRepo.Store(game); err != nil {
		return nil, err
	}

	return game, nil
}
