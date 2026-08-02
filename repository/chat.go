package repository

import (
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var chatIdentityColumns = []string{"title", "username"}

type ChatRepository struct {
	orm *gorm.DB
}

var _ domain.ChatRepository = (*ChatRepository)(nil)

func NewChatRepository(
	orm *gorm.DB,
) domain.ChatRepository {
	return &ChatRepository{
		orm: orm,
	}
}

func (repo ChatRepository) Get(id int64) (domain.Chat, error) {
	var chat domain.Chat
	if err := repo.orm.First(&chat, id).Error; err != nil {
		return chat, err
	}

	return chat, nil
}

func (repo ChatRepository) Save(chat *domain.Chat) error {
	return repo.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(chatIdentityColumns),
	}).Create(chat).Error
}

func (repo ChatRepository) UpdateSettings(chat *domain.Chat) error {
	return repo.orm.Updates(chat).Error
}
