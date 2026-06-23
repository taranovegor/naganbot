package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type DuelRepository struct {
	domain.DuelRepository
	orm *gorm.DB
}

func NewDuelRepository(orm *gorm.DB) domain.DuelRepository {
	return &DuelRepository{orm: orm}
}

func (repo DuelRepository) GetByID(id uuid.UUID) (*domain.Duel, error) {
	var duel domain.Duel
	if err := repo.orm.Preload("Challenger").Preload("Opponent").First(&duel, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &duel, nil
}

func (repo DuelRepository) Store(duel *domain.Duel) error {
	return repo.orm.Create(duel).Error
}

func (repo DuelRepository) Update(duel *domain.Duel) error {
	return repo.orm.Save(duel).Error
}
