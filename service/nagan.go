package service

import (
	"github.com/taranovegor/naganbot/domain"
	"math/rand"
)

type Nagan struct {
}

func NewNagan() *Nagan {
	return &Nagan{}
}

func (ng Nagan) Shot(gunslingers []*domain.Gunslinger) *domain.Gunslinger {
	return gunslingers[rand.Intn(len(gunslingers))]
}
