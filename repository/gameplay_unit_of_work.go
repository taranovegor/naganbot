package repository

import (
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type GameplayUnitOfWork struct {
	orm *gorm.DB
}

var _ domain.GameplayUnitOfWork = (*GameplayUnitOfWork)(nil)

func NewGameplayUnitOfWork(
	orm *gorm.DB,
) domain.GameplayUnitOfWork {
	return &GameplayUnitOfWork{
		orm: orm,
	}
}

func (repo GameplayUnitOfWork) CommitPlayedGame(game *domain.Game, victims []*domain.Gunslinger) error {
	return repo.orm.Transaction(func(tx *gorm.DB) error {
		if err := NewGameRepository(tx).Update(game); err != nil {
			return err
		}

		return NewGunslingerRepository(tx).Update(victims)
	})
}
