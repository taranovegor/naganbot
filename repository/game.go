package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
	"time"
)

type GameRepository struct {
	orm *gorm.DB
}

var _ domain.GameRepository = (*GameRepository)(nil)

func NewGameRepository(
	orm *gorm.DB,
) *GameRepository {
	return &GameRepository{
		orm: orm,
	}
}

func (repo GameRepository) GetByID(id uuid.UUID) (*domain.Game, error) {
	var game domain.Game
	err := repo.orm.Where("id = ?", id).First(&game).Error
	return &game, err
}

func (repo GameRepository) GetLatestGamesInChat(chatID int64, limit int) ([]domain.Game, error) {
	var games []domain.Game
	err := repo.getQueryByChat(chatID).
		Where("played_at IS NOT NULL").
		Order("played_at DESC").
		Limit(limit).
		Find(&games).
		Error

	return games, err
}

func (repo GameRepository) GetLatestForChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Order("created_at DESC").
		First(&game).
		Error
	return &game, err
}

func (repo GameRepository) GetActiveForChat(chatID int64, idleSince time.Time) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Scopes(activeSince(idleSince)).
		Order("created_at DESC").
		First(&game).
		Error
	return &game, err
}

func (repo GameRepository) Store(game *domain.Game) error {
	return repo.orm.Create(game).Error
}

func (repo GameRepository) Update(game *domain.Game) error {
	return repo.orm.Updates(game).Error
}

func (repo GameRepository) HasPlayedSince(chatID int64, since time.Time) bool {
	var counter int64
	repo.orm.Model(&domain.Game{}).
		Where("chat_id = ?", chatID).
		Where("played_at >= ?", since).
		Count(&counter)

	return counter > 0
}

// activeSince keeps unplayed games that were created or joined after idleSince.
func activeSince(idleSince time.Time) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("games.played_at IS NULL").
			Where(
				"games.created_at > ? OR EXISTS (SELECT 1 FROM gunslingers WHERE gunslingers.game_id = games.id AND gunslingers.joined_at > ?)",
				idleSince, idleSince,
			)
	}
}

func (repo GameRepository) getQueryByChat(chatID int64) *gorm.DB {
	return repo.orm.
		Preload("Gunslingers", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_at ASC")
		}).
		Preload("Gunslingers.Player").
		Preload("Owner").
		Where("chat_id = ?", chatID)
}
