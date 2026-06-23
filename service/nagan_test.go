package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
)

func createGunslingers(n int) []*domain.Gunslinger {
	gunslingers := make([]*domain.Gunslinger, n)
	for i := 0; i < n; i++ {
		gunslingers[i] = domain.NewGunslinger(uuid.New(), int64(i))
	}
	return gunslingers
}

func TestLeadBullet_Hit(t *testing.T) {
	b := NewLeadBullet()
	gunslingers := createGunslingers(5)
	result := b.Hit(gunslingers, 2)

	if len(result) != 1 {
		t.Fatalf("expected 1 gunslinger hit, got %d", len(result))
	}
	if result[0] != gunslingers[2] {
		t.Fatal("expected gunslinger at index 2")
	}
}

func TestAtomicBullet_Hit(t *testing.T) {
	b := NewAtomicBullet()
	gunslingers := createGunslingers(3)
	result := b.Hit(gunslingers, 0)

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
	defaultBullet := NewLeadBullet()
	factory := NewBulletFactory(defaultBullet)

	bullet := factory.Create(50)
	if bullet.Type() != BulletLeadType {
		t.Errorf("expected default bullet type %s, got %s", BulletLeadType, bullet.Type())
	}
}

func TestBulletFactory_Create_Special(t *testing.T) {
	defaultBullet := NewLeadBullet()
	specialBullet := NewAtomicBullet()

	factory := NewBulletFactory(defaultBullet, WeightedBullet{
		Chance: 100,
		Bullet: specialBullet,
	})

	bullet := factory.Create(0)
	if bullet.Type() != BulletAtomicType {
		t.Errorf("expected special bullet type %s, got %s", BulletAtomicType, bullet.Type())
	}
}

func TestNagan_Shoot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(drand.Beacon{
			Round:      12345,
			Randomness: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			Signature:  "sig",
		})
	}))
	defer server.Close()

	client := drand.NewClientWithURL(server.URL)
	gunslingers := createGunslingers(6)
	gameID := uuid.New()

	factory := NewBulletFactory(NewLeadBullet(), WeightedBullet{Chance: 3, Bullet: NewAtomicBullet()})
	nagan := NewNagan(factory, client)

	report, err := nagan.Shoot(context.Background(), gameID, gunslingers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.BulletType == "" {
		t.Error("expected bullet type to be set")
	}

	if len(report.Victims) == 0 {
		t.Error("expected at least one victim")
	}

	if report.ProofURL == "" {
		t.Error("expected proof URL to be set")
	}
}
