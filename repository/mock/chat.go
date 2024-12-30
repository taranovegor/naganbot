package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/domain"
)

type ChatRepository struct {
	mock.Mock
	domain.ChatRepository
}

func (m *ChatRepository) Exists(id int64) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *ChatRepository) Get(id int64) (domain.Chat, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Chat), args.Error(1)
}

func (m *ChatRepository) Store(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *ChatRepository) Update(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}
