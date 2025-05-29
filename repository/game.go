package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type GameRepository struct {
	domain.GameRepository
	orm *gorm.DB
}

func NewGameRepository(
	orm *gorm.DB,
) domain.GameRepository {
	return &GameRepository{
		orm: orm,
	}
}

func (repo GameRepository) GetByID(id uuid.UUID) (*domain.Game, error) {
	var game domain.Game
	err := repo.orm.Where("id = ?", id).First(&game).Error
	return &game, err
}

func (repo GameRepository) GetLatestForChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Order("created_at DESC").
		First(&game).
		Error
	return &game, err
}

func (repo GameRepository) GetActiveInChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Where("played_at IS NULL").
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

func (repo GameRepository) HasPendingInChat(chatID int64) bool {
	var counter int64
	repo.orm.Model(&domain.Game{}).
		Where("chat_id = ? AND played_at IS NULL", chatID).
		Count(&counter)

	return counter > 0
}

func (repo GameRepository) GetPrevInChatBeforeID(id uuid.UUID) (*domain.Game, error) {
	current, err := repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var prev domain.Game
	err = repo.getQueryByChat(current.ChatID).
		Where("created_at < ?", current.CreatedAt).
		Order("created_at DESC").
		First(&prev).Error

	if err != nil {
		return nil, err
	}

	return &prev, nil
}

func (repo GameRepository) GetNextInChatAfterID(id uuid.UUID) (*domain.Game, error) {
	current, err := repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var prev domain.Game
	err = repo.getQueryByChat(current.ChatID).
		Where("created_at > ?", current.CreatedAt).
		Order("created_at ASC").
		First(&prev).Error

	if err != nil {
		return nil, err
	}

	return &prev, nil
}

func (repo GameRepository) getQueryByChat(chatID int64) *gorm.DB {
	return repo.orm.
		Preload("Gunslingers", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_at ASC")
		}).
		Preload("Gunslingers.Player").
		Preload("Gunslingers.Game").
		Where("chat_id = ?", chatID)
}
