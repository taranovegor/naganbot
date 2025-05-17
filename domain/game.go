package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Game struct {
	ID          uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	ChatID      int64
	Chat        Chat
	OwnerID     int64
	Owner       User
	Gunslingers []*Gunslinger
	CreatedAt   time.Time
	PlayedAt    sql.NullTime
	BulletType  string
}

type GameRepository interface {
	GetByID(uuid.UUID) (*Game, error)
	GetLatestForChat(int64) (*Game, error)
	GetActiveForChat(int64) (*Game, error)
	Store(*Game) error
	Update(*Game) error
	HasActiveOrCreatedTodayInChat(id int64) bool
}

func NewGame(chatID int64, ownerID int64) *Game {
	ID := uuid.New()
	gunslinger := NewGunslinger(ID, ownerID)

	game := &Game{
		ID:          ID,
		ChatID:      chatID,
		OwnerID:     ownerID,
		Gunslingers: []*Gunslinger{gunslinger},
		CreatedAt:   time.Now(),
		PlayedAt:    sql.NullTime{Time: time.Now()},
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
