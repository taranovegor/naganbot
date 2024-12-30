package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/domain"
)

type UserRepository struct {
	mock.Mock
	domain.UserRepository
}

func (m *UserRepository) Exists(id int64) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *UserRepository) Get(id int64) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) Store(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
