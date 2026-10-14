//go:build integration

package repository

import (
	"errors"
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

type gameFixture struct {
	tx          *gorm.DB
	gameRepo    domain.GameRepository
	gunslingers domain.GunslingerRepository
	chatID      int64
	owner       *domain.User
}

func newGameFixture(t *testing.T) gameFixture {
	t.Helper()

	tx := withTx(t)
	f := gameFixture{
		tx:          tx,
		gameRepo:    NewGameRepository(tx),
		gunslingers: NewGunslingerRepository(tx),
		chatID:      testID(0),
		owner:       domain.NewUser(testID(1), "Owner", "", "owner"),
	}

	if err := NewChatRepository(tx).Save(domain.NewChat(f.chatID, "Chat", "chat")); err != nil {
		t.Fatalf("save chat: %v", err)
	}
	if err := NewUserRepository(tx).Save(f.owner); err != nil {
		t.Fatalf("save owner: %v", err)
	}

	return f
}

func (f gameFixture) newGame(t *testing.T, createdAt time.Time) *domain.Game {
	t.Helper()

	game := domain.NewGame(f.chatID, f.owner.ID, *f.owner, 6)
	game.CreatedAt = createdAt
	if err := f.gameRepo.Store(game); err != nil {
		t.Fatalf("store game: %v", err)
	}

	return game
}

func (f gameFixture) join(t *testing.T, game *domain.Game, joinedAt time.Time) {
	t.Helper()

	user := domain.NewUser(testID(100+int64(len(game.Gunslingers))+time.Now().UnixNano()%1000), "Player", "", "player")
	if err := NewUserRepository(f.tx).Save(user); err != nil {
		t.Fatalf("save player: %v", err)
	}

	gunslinger := domain.NewGunslinger(game.ID, *user)
	gunslinger.JoinedAt = joinedAt
	if err := f.gunslingers.Store(gunslinger); err != nil {
		t.Fatalf("store gunslinger: %v", err)
	}
}

func (f gameFixture) markPlayed(t *testing.T, game *domain.Game, playedAt time.Time) {
	t.Helper()

	if err := f.tx.Model(&domain.Game{}).Where("id = ?", game.ID).
		Update("played_at", playedAt).Error; err != nil {
		t.Fatalf("mark played: %v", err)
	}
}

func TestGameRepositoryGetActiveForChat(t *testing.T) {
	f := newGameFixture(t)
	idleSince := time.Now().Add(-domain.LobbyIdleTTL)

	if _, err := f.gameRepo.GetActiveForChat(f.chatID, idleSince); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected no active game in an empty chat, got %v", err)
	}

	stale := f.newGame(t, time.Now().Add(-2*domain.LobbyIdleTTL))
	f.join(t, stale, time.Now().Add(-2*domain.LobbyIdleTTL))
	if _, err := f.gameRepo.GetActiveForChat(f.chatID, idleSince); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected a lobby idle for longer than the TTL to be skipped, got %v", err)
	}

	f.join(t, stale, time.Now().Add(-time.Hour))
	got, err := f.gameRepo.GetActiveForChat(f.chatID, idleSince)
	if err != nil {
		t.Fatalf("expected a recent join to revive the old lobby, got %v", err)
	}
	if got.ID != stale.ID {
		t.Fatalf("expected the revived lobby, got %s", got.ID)
	}

	fresh := f.newGame(t, time.Now())
	got, err = f.gameRepo.GetActiveForChat(f.chatID, idleSince)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != fresh.ID {
		t.Fatalf("expected the newest active game to win, got %s", got.ID)
	}

	f.markPlayed(t, fresh, time.Now())
	f.markPlayed(t, stale, time.Now())
	if _, err := f.gameRepo.GetActiveForChat(f.chatID, idleSince); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected played games not to be active, got %v", err)
	}
}

func TestGameRepositoryHasPlayedSince(t *testing.T) {
	f := newGameFixture(t)
	midnight := time.Now().Truncate(24 * time.Hour)

	if f.gameRepo.HasPlayedSince(f.chatID, midnight) {
		t.Fatal("expected nothing played in an empty chat")
	}

	game := f.newGame(t, time.Now())
	if f.gameRepo.HasPlayedSince(f.chatID, midnight) {
		t.Fatal("expected an unplayed lobby not to count as played")
	}

	f.markPlayed(t, game, time.Now())
	if !f.gameRepo.HasPlayedSince(f.chatID, midnight) {
		t.Fatal("expected a game played after the cutoff to count")
	}
	if f.gameRepo.HasPlayedSince(f.chatID+1, midnight) {
		t.Fatal("expected another chat's games not to count")
	}

	f.markPlayed(t, game, time.Now().Add(-48*time.Hour))
	if f.gameRepo.HasPlayedSince(f.chatID, midnight) {
		t.Fatal("expected a game played before the cutoff not to count")
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
