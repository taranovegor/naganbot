package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Game struct {
	ID           uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	ChatID       int64
	Chat         Chat
	OwnerID      int64
	Owner        User
	Gunslingers  []*Gunslinger
	CreatedAt    time.Time
	PlayedAt     sql.NullTime
	BulletType   string
	PlayersCount int `gorm:"->;<-:create;not null;default:6"`
}

type GameRepository interface {
	GetByID(uuid.UUID) (*Game, error)
	GetLatestForChat(int64) (*Game, error)
	GetActiveInChat(int64) (*Game, error)
	Store(*Game) error
	Update(*Game) error
	HasPendingInChat(chatID int64) bool
	GetPrevInChatBeforeID(uuid.UUID) (*Game, error)
	GetNextInChatAfterID(id uuid.UUID) (*Game, error)
}

func NewGame(chatID int64, ownerID int64, playersCount int) *Game {
	ID := uuid.New()
	gunslinger := NewGunslinger(ID, ownerID)

	game := &Game{
		ID:           ID,
		ChatID:       chatID,
		OwnerID:      ownerID,
		CreatedAt:    time.Now(),
		PlayedAt:     sql.NullTime{Time: time.Now()},
		PlayersCount: playersCount,
	}
	gunslinger.Game = game

	return game
}

func (g *Game) IsPlayed() bool {
	return g.PlayedAt.Valid == true
}

func (g *Game) MarkAsPlayed(withBullet string) {
	if !g.IsPlayed() {
		g.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}
		g.BulletType = withBullet
	}
}
