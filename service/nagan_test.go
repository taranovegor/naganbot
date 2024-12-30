package service

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"testing"
)

func newTestNagan() Nagan {
	return NewNagan()
}

func TestNagan_Shot(t *testing.T) {
	nagan := newTestNagan()

	gameID := uuid.New()
	gunslingerFirst := domain.NewGunslinger(gameID, 1)
	gunslingerSecond := domain.NewGunslinger(gameID, 2)
	gunslingers := []*domain.Gunslinger{gunslingerFirst, gunslingerSecond}

	shooter := nagan.Shot(gunslingers)
	if gunslingerFirst.PlayerID != shooter.PlayerID && gunslingerSecond.PlayerID != shooter.PlayerID {
		t.Error()
	}
}
