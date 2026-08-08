//go:build integration

package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/taranovegor/naganbot/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		fmt.Println("skipping repository integration tests: TEST_DATABASE_DSN is not set")
		os.Exit(0)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&domain.Chat{}, &domain.User{}, &domain.Game{}, &domain.Gunslinger{}); err != nil {
		panic(err)
	}

	testDB = db
	os.Exit(m.Run())
}

func withTx(t *testing.T) *gorm.DB {
	t.Helper()
	tx := testDB.Begin()
	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil {
			t.Logf("rollback failed: %v", err)
		}
	})
	return tx
}

func testID(offset int64) int64 {
	return time.Now().UnixNano() + offset
}

func TestChatRepositorySavePreservesSettingsOnConflict(t *testing.T) {
	tx := withTx(t)
	repo := NewChatRepository(tx)

	chatID := testID(0)
	if err := repo.Save(domain.NewChat(chatID, "Original Title", "original_user")); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	saved, err := repo.Get(chatID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	saved.Settings.RequiredPlayers = 9
	if err := repo.UpdateSettings(&saved); err != nil {
		t.Fatalf("update settings failed: %v", err)
	}

	if err := repo.Save(domain.NewChat(chatID, "New Title", "new_user")); err != nil {
		t.Fatalf("second save failed: %v", err)
	}

	got, err := repo.Get(chatID)
	if err != nil {
		t.Fatalf("final get failed: %v", err)
	}
	if got.Title.String != "New Title" || got.Username.String != "new_user" {
		t.Fatalf("expected identity to be updated, got title=%q username=%q", got.Title.String, got.Username.String)
	}
	if got.Settings.RequiredPlayers != 9 {
		t.Fatalf("expected RequiredPlayers to survive the identity upsert, got %d", got.Settings.RequiredPlayers)
	}
}

func TestGameRepositoryHasActiveOrCreatedTodayInChat(t *testing.T) {
	tx := withTx(t)
	chatRepo := NewChatRepository(tx)
	userRepo := NewUserRepository(tx)
	gameRepo := NewGameRepository(tx)

	chatID := testID(0)
	owner := domain.NewUser(testID(1), "Owner", "", "owner")

	if err := chatRepo.Save(domain.NewChat(chatID, "Chat", "chat")); err != nil {
		t.Fatalf("save chat: %v", err)
	}
	if err := userRepo.Save(owner); err != nil {
		t.Fatalf("save owner: %v", err)
	}

	if gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		t.Fatal("expected no active or today's game before any game exists")
	}

	game := domain.NewGame(chatID, owner.ID, *owner, 6)
	if err := gameRepo.Store(game); err != nil {
		t.Fatalf("store game: %v", err)
	}

	if !gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		t.Fatal("expected a freshly created, unplayed game to count")
	}

	game.MarkAsPlayed("lead", "")
	if err := tx.Model(&domain.Game{}).Where("id = ?", game.ID).
		Update("played_at", game.PlayedAt).Error; err != nil {
		t.Fatalf("mark played: %v", err)
	}

	if !gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		t.Fatal("expected a game played earlier today to still block a new one (same-day cooldown)")
	}

	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	if err := tx.Model(&domain.Game{}).Where("id = ?", game.ID).
		Updates(map[string]interface{}{"created_at": twoDaysAgo, "played_at": twoDaysAgo}).Error; err != nil {
		t.Fatalf("backdate game: %v", err)
	}

	if gameRepo.HasActiveOrCreatedTodayInChat(chatID) {
		t.Fatal("expected a game created and played two days ago to no longer block a new game")
	}
}

func TestGunslingerRepositoryGetTopShotPlayersInChat(t *testing.T) {
	tx := withTx(t)
	chatRepo := NewChatRepository(tx)
	userRepo := NewUserRepository(tx)
	gameRepo := NewGameRepository(tx)
	gunslingerRepo := NewGunslingerRepository(tx)

	chatID := testID(0)
	loser := domain.NewUser(testID(1), "Loser", "", "loser")
	winner := domain.NewUser(testID(2), "Winner", "", "winner")

	if err := chatRepo.Save(domain.NewChat(chatID, "Chat", "chat")); err != nil {
		t.Fatalf("save chat: %v", err)
	}
	if err := userRepo.Save(loser); err != nil {
		t.Fatalf("save loser: %v", err)
	}
	if err := userRepo.Save(winner); err != nil {
		t.Fatalf("save winner: %v", err)
	}

	for i := 0; i < 2; i++ {
		game := domain.NewGame(chatID, loser.ID, *loser, 2)
		if err := gameRepo.Store(game); err != nil {
			t.Fatalf("store game: %v", err)
		}

		loserGunslinger := domain.NewGunslinger(game.ID, *loser)
		loserGunslinger.Game = game
		loserGunslinger.MarkAsShotHimself()
		if err := gunslingerRepo.Store(loserGunslinger); err != nil {
			t.Fatalf("store loser gunslinger: %v", err)
		}

		winnerGunslinger := domain.NewGunslinger(game.ID, *winner)
		if err := gunslingerRepo.Store(winnerGunslinger); err != nil {
			t.Fatalf("store winner gunslinger: %v", err)
		}

		game.MarkAsPlayed("lead", "")
		if err := tx.Model(&domain.Game{}).Where("id = ?", game.ID).
			Update("played_at", game.PlayedAt).Error; err != nil {
			t.Fatalf("mark played: %v", err)
		}
	}

	players, err := gunslingerRepo.GetTopShotPlayersInChat(chatID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(players) != 1 {
		t.Fatalf("expected exactly one player with self-shots, got %d: %+v", len(players), players)
	}
	if players[0].PlayerId != loser.ID || players[0].Times != 2 {
		t.Fatalf("expected loser %d with 2 self-shots, got %+v", loser.ID, players[0])
	}
}
