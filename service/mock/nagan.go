package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

type Nagan struct {
	mock.Mock
	service.Nagan
}

func (m *Nagan) Shot(gunslingers []*domain.Gunslinger) *domain.Gunslinger {
	args := m.Called(gunslingers)
	return args.Get(0).(*domain.Gunslinger)
}
