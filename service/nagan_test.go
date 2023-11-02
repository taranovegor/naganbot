package service

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"math/rand"
	"testing"
)

func createGunslingers(n int) []*domain.Gunslinger {
	gunslingers := make([]*domain.Gunslinger, n)
	for i := 0; i < n; i++ {
		gunslingers[i] = domain.NewGunslinger(uuid.New(), int64(i))
	}
	return gunslingers
}

func TestLeadBullet_Hit(t *testing.T) {
	rand.Seed(1)
	b := NewLeadBullet()
	gunslingers := createGunslingers(5)
	result := b.Hit(gunslingers)

	if len(result) != 1 {
		t.Fatalf("expected 1 gunslinger hit, got %d", len(result))
	}
	if result[0] == nil {
		t.Fatal("expected valid gunslinger, got nil")
	}
}

func TestAtomicBullet_Hit(t *testing.T) {
	b := NewAtomicBullet()
	gunslingers := createGunslingers(3)
	result := b.Hit(gunslingers)

	if len(result) != len(gunslingers) {
		t.Fatalf("expected %d gunslingers hit, got %d", len(gunslingers), len(result))
	}
	for i, g := range result {
		if g != gunslingers[i] {
			t.Errorf("expected gunslinger %v at index %d, got %v", gunslingers[i], i, g)
		}
	}
}

func TestBulletFactory_Create_Default(t *testing.T) {
	rand.Seed(99)
	defaultBullet := NewLeadBullet()
	factory := NewBulletFactory(defaultBullet)

	bullet := factory.Create()
	if bullet.Type() != BulletLeadType {
		t.Errorf("expected default bullet type %s, got %s", BulletLeadType, bullet.Type())
	}
}

func TestBulletFactory_Create_Special(t *testing.T) {
	rand.Seed(1)
	defaultBullet := NewLeadBullet()
	specialBullet := NewAtomicBullet()

	factory := NewBulletFactory(defaultBullet, WeightedBullet{
		Chance: 100,
		Bullet: specialBullet,
	})

	bullet := factory.Create()
	if bullet.Type() != BulletAtomicType {
		t.Errorf("expected special bullet type %s, got %s", BulletAtomicType, bullet.Type())
	}
}

func TestNagan_Shoot(t *testing.T) {
	rand.Seed(1)
	gunslingers := createGunslingers(3)
	defaultBullet := NewLeadBullet()
	specialBullet := NewAtomicBullet()
	factory := NewBulletFactory(defaultBullet, WeightedBullet{
		Chance: 100,
		Bullet: specialBullet,
	})
	nagan := NewNagan(factory)

	report := nagan.Shoot(gunslingers)

	if report.BulletType != BulletAtomicType {
		t.Errorf("expected bullet type %s, got %s", BulletAtomicType, report.BulletType)
	}

	if len(report.Gunslingers) != 3 {
		t.Errorf("expected 3 gunslingers to be hit, got %d", len(report.Gunslingers))
	}
}
