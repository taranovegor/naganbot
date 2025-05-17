package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type GunslingerRepository struct {
	domain.GunslingerRepository
	orm *gorm.DB
}

func NewGunslingerRepository(
	orm *gorm.DB,
) domain.GunslingerRepository {
	return &GunslingerRepository{
		orm: orm,
	}
}

func (repo GunslingerRepository) Store(gunslinger *domain.Gunslinger) error {
	return repo.orm.Create(gunslinger).Error
}

// Update todo: optimize query
func (repo GunslingerRepository) Update(gunslingers []*domain.Gunslinger) error {
	orm := repo.orm
	for _, gunslinger := range gunslingers {
		orm.Updates(gunslinger)
	}
	return orm.Error
}

func (repo GunslingerRepository) GetByGameID(gameID uuid.UUID) ([]*domain.Gunslinger, error) {
	var gunslingers []*domain.Gunslinger
	err := repo.orm.
		Preload("Game").
		Preload("Player").
		Where("game_id = ?", gameID).
		Order("joined_at ASC").
		Find(&gunslingers).
		Error

	return gunslingers, err
}

func (repo GunslingerRepository) IsPlayerExistsInGame(userID int64, gameID uuid.UUID) bool {
	var counter int64
	repo.orm.Model(&domain.Gunslinger{}).
		Where("player_id = ?", userID).
		Where("game_id = ?", gameID).
		Count(&counter)

	return counter > 0
}

func (repo GunslingerRepository) GetTopShotPlayersInChat(chatID int64) ([]domain.GunslingerTopShotPlayer, error) {
	var players []domain.GunslingerTopShotPlayer
	err := repo.getQueryTopShotPlayersInChat(chatID).
		Find(&players).
		Error

	return players, err
}

func (repo GunslingerRepository) GetTopShopPlayersByYearInChat(chatID int64, year int) ([]domain.GunslingerTopShotPlayer, error) {
	var players []domain.GunslingerTopShotPlayer
	err := repo.getQueryTopShotPlayersInChat(chatID).
		Where("YEAR(Game.created_at) = ?", year).
		Find(&players).
		Error

	return players, err
}

func (repo GunslingerRepository) CountNumberOfPlayerGamesInChat(userID int64, chatID int64) int64 {
	var counter int64
	repo.getQueryPlayerGamesInChat(userID, chatID).
		Count(&counter)

	return counter
}

func (repo GunslingerRepository) CountNumberOfSelfShotsInChat(userID int64, chatID int64) int64 {
	var counter int64
	repo.getQueryPlayerGamesInChat(userID, chatID).
		Where("shot_himself = ?", true).
		Count(&counter)

	return counter
}

func (repo GunslingerRepository) getQueryTopShotPlayersInChat(chatID int64) *gorm.DB {
	return repo.orm.Model(&domain.Gunslinger{}).
		Select("player_id, COUNT(chat_id) as times").
		InnerJoins("Game").
		Where("chat_id = ?", chatID).
		Where("played_at IS NOT NULL").
		Where("shot_himself = ?", true).
		Group("player_id").
		Order("times DESC").
		Limit(10)
}

func (repo GunslingerRepository) getQueryPlayerGamesInChat(userID int64, chatID int64) *gorm.DB {
	return repo.orm.Model(&domain.Gunslinger{}).
		InnerJoins("Game").
		Where("player_id = ?", userID).
		Where("chat_id = ?", chatID).
		Where("played_at IS NOT NULL")
}
