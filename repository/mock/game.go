package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type GameRepository struct {
	mock.Mock
	domain.GameRepository
}

func (m *GameRepository) GetLatestForChat(chatID int64) (domain.Game, error) {
	args := m.Called(chatID)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *GameRepository) GetActiveForChat(chatID int64) (domain.Game, error) {
	args := m.Called(chatID)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *GameRepository) Store(game *domain.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *GameRepository) Update(game *domain.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *GameRepository) HasActiveOrCreatedTodayInChat(chatID int64) bool {
	args := m.Called(chatID)
	return args.Bool(0)
}

func (m *GameRepository) getQueryByChat(chatID int64) *gorm.DB {
	args := m.Called(chatID)
	return args.Get(0).(*gorm.DB)
}
