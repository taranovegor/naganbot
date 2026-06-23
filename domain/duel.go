package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Duel struct {
	ID           uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	ChatID       int64
	ChallengerID int64
	Challenger   User
	OpponentID   int64
	Opponent     User
	CreatedAt    time.Time
	PlayedAt     sql.NullTime
	LoserID      int64
	ProofURL     string
}

type DuelRepository interface {
	GetByID(uuid.UUID) (*Duel, error)
	Store(*Duel) error
	Update(*Duel) error
}

func NewDuel(chatID int64, challenger User, opponent User) *Duel {
	return &Duel{
		ID:           uuid.Must(uuid.NewV7()),
		ChatID:       chatID,
		ChallengerID: challenger.ID,
		Challenger:   challenger,
		OpponentID:   opponent.ID,
		Opponent:     opponent,
		CreatedAt:    time.Now(),
	}
}

func (d *Duel) IsPlayed() bool {
	return d.PlayedAt.Valid
}

func (d *Duel) MarkAsPlayed(loserID int64, proofURL string) {
	if !d.IsPlayed() {
		d.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}
		d.LoserID = loserID
		d.ProofURL = proofURL
	}
}
