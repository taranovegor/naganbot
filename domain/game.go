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
	ProofURL     string
	PlayersCount int `gorm:"->;<-:create;not null;default:6"`
}

type GameRepository interface {
	GetByID(uuid.UUID) (*Game, error)
	GetLatestGamesInChat(int64, int) ([]Game, error)
	GetLatestForChat(int64) (*Game, error)
	GetActiveForChat(int64) (*Game, error)
	Store(*Game) error
	HasActiveOrCreatedTodayInChat(id int64) bool
}

func NewGame(chatID int64, ownerID int64, owner User, playersCount int) *Game {
	return &Game{
		ID:           uuid.Must(uuid.NewV7()),
		ChatID:       chatID,
		OwnerID:      ownerID,
		Owner:        owner,
		CreatedAt:    time.Now(),
		PlayedAt:     sql.NullTime{},
		PlayersCount: playersCount,
	}
}

func (g *Game) IsPlayed() bool {
	return g.PlayedAt.Valid
}

func (g *Game) MarkAsPlayed(withBullet string, proofURL string) {
	if !g.IsPlayed() {
		g.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}
		g.BulletType = withBullet
		g.ProofURL = proofURL
	}
}
