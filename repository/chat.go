package repository

import (
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type ChatRepository struct {
	domain.ChatRepository
	orm *gorm.DB
}

func NewChatRepository(
	orm *gorm.DB,
) domain.ChatRepository {
	return &ChatRepository{
		orm: orm,
	}
}

func (repo ChatRepository) Exists(id int64) bool {
	var counter int64
	repo.orm.Model(&domain.Chat{}).Where(id).Count(&counter)

	return counter > 0
}

func (repo ChatRepository) Get(id int64) (domain.Chat, error) {
	var chat domain.Chat
	if err := repo.orm.First(&chat, id).Error; err != nil {
		return chat, err
	}

	return chat, nil
}

func (repo ChatRepository) Store(chat *domain.Chat) error {
	return repo.orm.Create(chat).Error
}

func (repo ChatRepository) Update(chat *domain.Chat) error {
	return repo.orm.Updates(chat).Error
}
