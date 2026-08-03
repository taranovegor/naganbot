package usecase

import (
	"errors"
	"fmt"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"gorm.io/gorm"
)

var (
	ErrGameCooldown = errors.New("game cooldown")
)

type CreateGameUseCase struct {
	locker   service.Locker
	chatRepo domain.ChatRepository
	gameRepo domain.GameRepository
	userRepo domain.UserRepository
}

func NewCreateGameUseCase(
	locker service.Locker,
	chatRepo domain.ChatRepository,
	gameRepo domain.GameRepository,
	userRepo domain.UserRepository,
) *CreateGameUseCase {
	return &CreateGameUseCase{
		locker:   locker,
		chatRepo: chatRepo,
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

func (uc *CreateGameUseCase) Execute(chatID int64, ownerID int64) (*domain.Game, error) {
	release, ok := uc.locker.TryLock(fmt.Sprintf("game-start-%d", chatID))
	if !ok {
		return nil, service.ErrLockNotAcquired
	}
	defer release()

	if uc.gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		game, err := uc.gameRepo.GetActiveForChat(chatID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrGameCooldown
			}

			return nil, err
		}

		return game, nil
	}

	owner, err := uc.userRepo.Get(ownerID)
	if err != nil {
		return nil, err
	}

	chat, err := uc.chatRepo.Get(chatID)
	if err != nil {
		return nil, err
	}

	game := domain.NewGame(chatID, ownerID, owner, chat.Settings.RequiredPlayers)

	if err := uc.gameRepo.Store(game); err != nil {
		return nil, err
	}

	return game, nil
}
