package mock

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/domain"
)

type GunslingerRepository struct {
	mock.Mock
	domain.GunslingerRepository
}

func (m *GunslingerRepository) Store(gunslinger *domain.Gunslinger) error {
	args := m.Called(gunslinger)
	return args.Error(0)
}

func (m *GunslingerRepository) Update(gunslinger *domain.Gunslinger) error {
	args := m.Called(gunslinger)
	return args.Error(0)
}

func (m *GunslingerRepository) IsPlayerExistsInGame(userID int64, gameID uuid.UUID) bool {
	args := m.Called(userID, gameID)
	return args.Bool(0)
}

func (m *GunslingerRepository) GetTopShotPlayersInChat(chatID int64) ([]domain.GunslingerTopShotPlayer, error) {
	args := m.Called(chatID)
	return args.Get(0).([]domain.GunslingerTopShotPlayer), args.Error(1)
}

func (m *GunslingerRepository) GetTopShotPlayersByYearInChat(chatID int64, year int) ([]domain.GunslingerTopShotPlayer, error) {
	args := m.Called(chatID, year)
	return args.Get(0).([]domain.GunslingerTopShotPlayer), args.Error(1)
}

func (m *GunslingerRepository) CountNumberOfPlayerGamesInChat(userID int64, chatID int64) int64 {
	args := m.Called(userID, chatID)
	return args.Get(0).(int64)
}

func (m *GunslingerRepository) CountNumberOfSelfShotsInChat(userID int64, chatID int64) int64 {
	args := m.Called(userID, chatID)
	return args.Get(0).(int64)
}
