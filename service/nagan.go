package service

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
)

const (
	BulletLeadType   string = "lead"
	BulletAtomicType string = "atomic"
)

type Bullet interface {
	Type() string
	Hit(gunslingers []*domain.Gunslinger, index int) []*domain.Gunslinger
}

type leadBullet struct{}

func NewLeadBullet() Bullet {
	return &leadBullet{}
}

func (b *leadBullet) Type() string {
	return BulletLeadType
}

func (b *leadBullet) Hit(gunslingers []*domain.Gunslinger, index int) []*domain.Gunslinger {
	return []*domain.Gunslinger{gunslingers[index]}
}

type atomicBullet struct{}

func NewAtomicBullet() Bullet {
	return &atomicBullet{}
}

func (b *atomicBullet) Type() string {
	return BulletAtomicType
}

func (b *atomicBullet) Hit(gunslingers []*domain.Gunslinger, _ int) []*domain.Gunslinger {
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

func (f *BulletFactory) Create(roll int) Bullet {
	acc := 0
	for _, wb := range f.specialBullets {
		acc += wb.Chance
		if roll < acc {
			return wb.Bullet
		}
	}
	return f.defaultBullet
}

type HitReport struct {
	Victims    []*domain.Gunslinger
	BulletType string
	ProofURL   string
}

type Nagan struct {
	bulletFactory *BulletFactory
	drand         *drand.Client
}

func NewNagan(bulletFactory *BulletFactory, drand *drand.Client) *Nagan {
	return &Nagan{
		bulletFactory: bulletFactory,
		drand:         drand,
	}
}

func (ng *Nagan) Shoot(ctx context.Context, gameID uuid.UUID, gunslingers []*domain.Gunslinger) (*HitReport, error) {
	beacon, err := ng.drand.GetLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get drand beacon: %w", err)
	}

	randomnessBytes, err := hex.DecodeString(beacon.Randomness)
	if err != nil {
		return nil, fmt.Errorf("failed to decode drand randomness: %w", err)
	}

	hash := sha256.Sum256(append(randomnessBytes, []byte(gameID.String())...))

	bulletRoll := int(binary.BigEndian.Uint32(hash[:4]) % 100)
	bullet := ng.bulletFactory.Create(bulletRoll)

	victimIndex := int(binary.BigEndian.Uint64(hash[4:12]) % uint64(len(gunslingers)))
	victims := bullet.Hit(gunslingers, victimIndex)

	return &HitReport{
		Victims:    victims,
		BulletType: bullet.Type(),
		ProofURL:   ng.drand.ProofURL(beacon.Round, gameID),
	}, nil
}
