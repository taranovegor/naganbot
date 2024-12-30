package service

import (
	"github.com/taranovegor/naganbot/domain"
	"math/rand"
)

type Nagan interface {
	Shot(gunslingers []*domain.Gunslinger) *domain.Gunslinger
}

type nagan struct {
	Nagan
}

func NewNagan() Nagan {
	return &nagan{}
}

func (ng *nagan) Shot(gunslingers []*domain.Gunslinger) *domain.Gunslinger {
	return gunslingers[rand.Intn(len(gunslingers))]
}
