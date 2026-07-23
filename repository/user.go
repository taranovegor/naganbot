package repository

import (
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	domain.UserRepository
	orm *gorm.DB
}

func NewUserRepository(
	orm *gorm.DB,
) domain.UserRepository {
	return &UserRepository{
		orm: orm,
	}
}

func (repo UserRepository) Get(id int64) (domain.User, error) {
	var user domain.User
	if err := repo.orm.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repo UserRepository) GetByIDs(ids []int64) ([]domain.User, error) {
	var users []domain.User
	if err := repo.orm.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (repo UserRepository) Save(user *domain.User) error {
	return repo.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(user).Error
}
