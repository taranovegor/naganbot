package usecase

import (
	"errors"
	"fmt"
	"time"

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

	now := time.Now()

	game, err := uc.gameRepo.GetActiveForChat(chatID, now.Add(-domain.LobbyIdleTTL))
	if err == nil {
		return game, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if uc.gameRepo.HasPlayedSince(chatID, startOfDay(now)) {
		return nil, ErrGameCooldown
	}

	owner, err := uc.userRepo.Get(ownerID)
	if err != nil {
		return nil, err
	}

	chat, err := uc.chatRepo.Get(chatID)
	if err != nil {
		return nil, err
	}

	game = domain.NewGame(chatID, ownerID, owner, chat.Settings.RequiredPlayers, chat.Settings.Mode)

	if err := uc.gameRepo.Store(game); err != nil {
		return nil, err
	}

	return game, nil
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
