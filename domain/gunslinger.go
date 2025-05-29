package domain

import (
	"github.com/google/uuid"
	"time"
)

type Gunslinger struct {
	ID          uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	GameID      uuid.UUID
	Game        *Game
	PlayerID    int64
	Player      User
	JoinedAt    time.Time
	ShotHimself bool
}

type GunslingerTopShotPlayer struct {
	PlayerId int64
	Times    int
}

type GunslingerRepository interface {
	Store(*Gunslinger) error
	Update([]*Gunslinger) error
	GetByGameID(uuid.UUID) ([]*Gunslinger, error)
	IsUserInGame(userID int64, gameID uuid.UUID) bool
	GetTopShotPlayersInChat(int64) ([]GunslingerTopShotPlayer, error)
	GetTopShopPlayersByYearInChat(chatID int64, year int) ([]GunslingerTopShotPlayer, error)
	CountNumberOfPlayerGamesInChat(userID int64, chatID int64) int64
	CountNumberOfSelfShotsInChat(userID int64, chatID int64) int64
	HasPlayedToday(userID int64) bool
}

func NewGunslinger(gameID uuid.UUID, playerID int64) *Gunslinger {
	return &Gunslinger{
		ID:          uuid.New(),
		GameID:      gameID,
		PlayerID:    playerID,
		JoinedAt:    time.Now(),
		ShotHimself: false,
	}
}

func (g *Gunslinger) MarkAsShotHimself() {
	if !g.Game.IsPlayed() {
		g.ShotHimself = true
	}
}
