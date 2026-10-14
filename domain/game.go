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
	PlayersCount int      `gorm:"->;<-:create;not null;default:6"`
	Mode         GameMode `gorm:"->;<-:create;not null;default:'fixed'"`
}

type GameRepository interface {
	GetByID(uuid.UUID) (*Game, error)
	GetLatestGamesInChat(int64, int) ([]Game, error)
	GetLatestForChat(int64) (*Game, error)
	// GetActiveForChat returns the newest unplayed game in the chat. It skips
	// lobbies with no activity since idleSince.
	GetActiveForChat(chatID int64, idleSince time.Time) (*Game, error)
	// GetActiveCountdownGames returns unplayed countdown games with their
	// gunslingers. It skips games with no activity since idleSince.
	GetActiveCountdownGames(idleSince time.Time) ([]*Game, error)
	Store(*Game) error
	HasPlayedSince(chatID int64, since time.Time) bool
}

func NewGame(chatID int64, ownerID int64, owner User, playersCount int, mode GameMode) *Game {
	return &Game{
		ID:           uuid.Must(uuid.NewV7()),
		ChatID:       chatID,
		OwnerID:      ownerID,
		Owner:        owner,
		CreatedAt:    time.Now(),
		PlayedAt:     sql.NullTime{},
		PlayersCount: playersCount,
		Mode:         mode,
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
