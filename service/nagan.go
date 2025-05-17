package service

import (
	"github.com/taranovegor/naganbot/domain"
	"math/rand"
)

const (
	BulletLeadType   string = "lead"
	BulletAtomicType string = "atomic"
)

type Bullet interface {
	Type() string
	Hit([]*domain.Gunslinger) []*domain.Gunslinger
}

type leadBullet struct{}

func NewLeadBullet() Bullet {
	return &leadBullet{}
}

func (b *leadBullet) Type() string {
	return BulletLeadType
}

func (b *leadBullet) Hit(gunslingers []*domain.Gunslinger) []*domain.Gunslinger {
	return []*domain.Gunslinger{gunslingers[rand.Intn(len(gunslingers))]}
}

type atomicBullet struct{}

func NewAtomicBullet() Bullet {
	return &atomicBullet{}
}

func (b *atomicBullet) Type() string {
	return BulletAtomicType
}

func (b *atomicBullet) Hit(gunslingers []*domain.Gunslinger) []*domain.Gunslinger {
	return gunslingers
}

type WeightedBullet struct {
	Chance int
	Bullet Bullet
}

type BulletFactory struct {
	defaultBullet  Bullet
	specialBullets []WeightedBullet
}

func NewBulletFactory(
	defaultBullet Bullet,
	specialBullets ...WeightedBullet,
) *BulletFactory {
	return &BulletFactory{
		defaultBullet:  defaultBullet,
		specialBullets: specialBullets,
	}
}

func (f *BulletFactory) Create() Bullet {
	r := rand.Intn(100)
	acc := 0
	for _, wb := range f.specialBullets {
		acc += wb.Chance
		if r < acc {
			return wb.Bullet
		}
	}
	return f.defaultBullet
}

type HitReport struct {
	Victims    []*domain.Gunslinger
	BulletType string
}

type Nagan struct {
	bulletFactory *BulletFactory
}

func NewNagan(bulletFactory *BulletFactory) *Nagan {
	return &Nagan{
		bulletFactory: bulletFactory,
	}
}

func (ng *Nagan) Shoot(gunslingers []*domain.Gunslinger) *HitReport {
	bullet := ng.bulletFactory.Create()
	gunslingers = bullet.Hit(gunslingers)
	return &HitReport{
		Victims:    gunslingers,
		BulletType: bullet.Type(),
	}
}
