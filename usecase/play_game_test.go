package usecase

import (
	"testing"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/service"
)

func TestGameLockKeyIsUniquePerGame(t *testing.T) {
	first := uuid.Must(uuid.NewV7())
	second := uuid.Must(uuid.NewV7())

	locker := service.NewLocker()
	firstMutex := locker.LockFor(gameLockKey(first))
	secondMutex := locker.LockFor(gameLockKey(second))

	if firstMutex == secondMutex {
		t.Fatal("two distinct games were assigned the same lock")
	}
}
