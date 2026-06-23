package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// DuelTTL is how long the opponent has to answer.
const DuelTTL = 10 * time.Minute

type Duel struct {
	ID           uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	ChatID       int64
	ChallengerID int64
	Challenger   User
	OpponentID   int64
	Opponent     User
	CreatedAt    time.Time
	PlayedAt     sql.NullTime
	DeclinedAt   sql.NullTime
	LoserID      int64
	ProofURL     string
}

type DuelRepository interface {
	GetByID(uuid.UUID) (*Duel, error)
	Store(*Duel) error
	// Resolve stores the outcome only if the duel is still open and reports
	// whether this call closed it.
	Resolve(*Duel) (bool, error)
	HasPendingFrom(chatID int64, challengerID int64, since time.Time) (bool, error)
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

func (d *Duel) IsDeclined() bool {
	return d.DeclinedAt.Valid
}

func (d *Duel) IsClosed() bool {
	return d.IsPlayed() || d.IsDeclined()
}

func (d *Duel) IsExpired(now time.Time) bool {
	return now.Sub(d.CreatedAt) > DuelTTL
}

func (d *Duel) MarkAsDeclined() {
	if !d.IsClosed() {
		d.DeclinedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
}

func (d *Duel) MarkAsPlayed(loserID int64, proofURL string) {
	if !d.IsClosed() {
		d.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}
		d.LoserID = loserID
		d.ProofURL = proofURL
	}
}
