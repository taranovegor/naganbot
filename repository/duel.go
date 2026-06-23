package repository

import (
	"time"

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

func (repo DuelRepository) Resolve(duel *domain.Duel) (bool, error) {
	res := repo.orm.Model(&domain.Duel{}).
		Where("id = ? AND played_at IS NULL AND declined_at IS NULL", duel.ID).
		Updates(map[string]interface{}{
			"played_at":   duel.PlayedAt,
			"declined_at": duel.DeclinedAt,
			"loser_id":    duel.LoserID,
			"proof_url":   duel.ProofURL,
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (repo DuelRepository) HasPendingFrom(chatID int64, challengerID int64, since time.Time) (bool, error) {
	var count int64
	err := repo.orm.Model(&domain.Duel{}).
		Where("chat_id = ? AND challenger_id = ? AND created_at > ? AND played_at IS NULL AND declined_at IS NULL", chatID, challengerID, since).
		Count(&count).Error
	return count > 0, err
}
