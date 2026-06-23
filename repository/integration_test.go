//go:build integration

package repository

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
	if err := db.AutoMigrate(&domain.Chat{}, &domain.User{}, &domain.Game{}, &domain.Gunslinger{}, &domain.Duel{}); err != nil {
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

func (f gameFixture) newGame(t *testing.T, mode domain.GameMode, createdAt time.Time) *domain.Game {
	t.Helper()

	game := domain.NewGame(f.chatID, f.owner.ID, *f.owner, 6, mode)
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

func TestGameRepositoryStoresTheModeAndDefaultsExistingRowsToFixed(t *testing.T) {
	f := newGameFixture(t)

	countdown := f.newGame(t, domain.GameModeCountdown, time.Now())
	got, err := f.gameRepo.GetByID(countdown.ID)
	if err != nil {
		t.Fatalf("get game: %v", err)
	}
	if got.Mode != domain.GameModeCountdown {
		t.Fatalf("expected the countdown mode to round-trip, got %q", got.Mode)
	}

	if err := f.tx.Exec("INSERT INTO chats (id) VALUES (?)", f.chatID+1).Error; err != nil {
		t.Fatalf("insert bare chat: %v", err)
	}
	chat, err := NewChatRepository(f.tx).Get(f.chatID + 1)
	if err != nil {
		t.Fatalf("get chat: %v", err)
	}
	if chat.Settings.Mode != domain.GameModeFixed {
		t.Fatalf("expected a chat without an explicit mode to be fixed, got %q", chat.Settings.Mode)
	}
}

func TestGameRepositoryGetActiveForChat(t *testing.T) {
	f := newGameFixture(t)
	idleSince := time.Now().Add(-domain.LobbyIdleTTL)

	if _, err := f.gameRepo.GetActiveForChat(f.chatID, idleSince); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected no active game in an empty chat, got %v", err)
	}

	stale := f.newGame(t, domain.GameModeFixed, time.Now().Add(-2*domain.LobbyIdleTTL))
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

	fresh := f.newGame(t, domain.GameModeFixed, time.Now())
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

	game := f.newGame(t, domain.GameModeFixed, time.Now())
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

func TestGameRepositoryGetActiveCountdownGames(t *testing.T) {
	f := newGameFixture(t)
	idleSince := time.Now().Add(-domain.LobbyIdleTTL)

	countdown := f.newGame(t, domain.GameModeCountdown, time.Now())
	f.join(t, countdown, time.Now().Add(-2*time.Minute))
	f.join(t, countdown, time.Now().Add(-time.Minute))

	f.newGame(t, domain.GameModeFixed, time.Now())
	f.newGame(t, domain.GameModeCountdown, time.Now().Add(-2*domain.LobbyIdleTTL))

	played := f.newGame(t, domain.GameModeCountdown, time.Now())
	f.markPlayed(t, played, time.Now())

	games, err := f.gameRepo.GetActiveCountdownGames(idleSince)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ours []*domain.Game
	for _, game := range games {
		if game.ChatID == f.chatID {
			ours = append(ours, game)
		}
	}
	if len(ours) != 1 || ours[0].ID != countdown.ID {
		t.Fatalf("expected only the live countdown lobby, got %d games", len(ours))
	}
	if len(ours[0].Gunslingers) != 2 {
		t.Fatalf("expected the gunslingers to be preloaded, got %d", len(ours[0].Gunslingers))
	}
	if !ours[0].Gunslingers[0].JoinedAt.Before(ours[0].Gunslingers[1].JoinedAt) {
		t.Fatal("expected the gunslingers ordered by join time")
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
		game := domain.NewGame(chatID, loser.ID, *loser, 2, domain.GameModeFixed)
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

func TestUserRepositoryGetByUsername(t *testing.T) {
	tx := withTx(t)
	repo := NewUserRepository(tx)
	username := fmt.Sprintf("u%d", testID(0))

	if err := repo.Save(domain.NewUser(testID(1), "Name", "", username)); err != nil {
		t.Fatalf("save user: %v", err)
	}

	got, err := repo.GetByUsername(username)
	if err != nil || got.Username != username {
		t.Fatalf("expected the saved user, got %+v, err %v", got, err)
	}
	if upper, err := repo.GetByUsername(strings.ToUpper(username)); err != nil || upper.Username != username {
		t.Fatalf("expected a case-insensitive match, got %+v, err %v", upper, err)
	}
	if _, err := repo.GetByUsername(username + "-missing"); err == nil {
		t.Fatal("expected an error for an unknown username")
	}
}

func TestDuelRepositoryStoreGetAndResolve(t *testing.T) {
	tx := withTx(t)
	users := NewUserRepository(tx)
	repo := NewDuelRepository(tx)

	challenger := domain.NewUser(testID(1), "Challenger", "", "challenger")
	opponent := domain.NewUser(testID(2), "Opponent", "", "opponent")
	for _, u := range []*domain.User{challenger, opponent} {
		if err := users.Save(u); err != nil {
			t.Fatalf("save user: %v", err)
		}
	}

	duel := domain.NewDuel(testID(3), *challenger, *opponent)
	if err := repo.Store(duel); err != nil {
		t.Fatalf("store duel: %v", err)
	}

	got, err := repo.GetByID(duel.ID)
	if err != nil {
		t.Fatalf("get duel: %v", err)
	}
	if got.Challenger.ID != challenger.ID || got.Opponent.ID != opponent.ID {
		t.Fatalf("expected participants to be preloaded, got %+v", got)
	}
	if got.IsPlayed() {
		t.Fatal("a stored duel must not be played")
	}

	got.MarkAsPlayed(opponent.ID, "https://example.com/proof")
	if claimed, err := repo.Resolve(got); err != nil || !claimed {
		t.Fatalf("resolve duel: claimed=%v err=%v", claimed, err)
	}

	again, err := repo.GetByID(duel.ID)
	if err != nil {
		t.Fatalf("get duel: %v", err)
	}
	if !again.IsPlayed() || again.LoserID != opponent.ID || again.ProofURL != "https://example.com/proof" {
		t.Fatalf("update was not persisted: %+v", again)
	}
}

func TestDuelRepositoryResolveClosesTheDuelOnlyOnce(t *testing.T) {
	tx := withTx(t)
	repo := NewDuelRepository(tx)
	users := NewUserRepository(tx)

	challenger := domain.NewUser(testID(1), "Challenger", "", "challenger")
	opponent := domain.NewUser(testID(2), "Opponent", "", "opponent")
	for _, u := range []*domain.User{challenger, opponent} {
		if err := users.Save(u); err != nil {
			t.Fatalf("save user: %v", err)
		}
	}
	duel := domain.NewDuel(testID(3), *challenger, *opponent)
	if err := repo.Store(duel); err != nil {
		t.Fatalf("store duel: %v", err)
	}

	first := *duel
	first.MarkAsPlayed(challenger.ID, "first")
	if claimed, err := repo.Resolve(&first); err != nil || !claimed {
		t.Fatalf("first resolve: claimed=%v err=%v", claimed, err)
	}

	second := *duel
	second.MarkAsPlayed(opponent.ID, "second")
	if claimed, err := repo.Resolve(&second); err != nil || claimed {
		t.Fatalf("second resolve must not claim: claimed=%v err=%v", claimed, err)
	}

	got, err := repo.GetByID(duel.ID)
	if err != nil || got.LoserID != challenger.ID || got.ProofURL != "first" {
		t.Fatalf("the first outcome must win, got %+v err=%v", got, err)
	}
}

func TestDuelRepositoryResolveStoresDeclines(t *testing.T) {
	tx := withTx(t)
	repo := NewDuelRepository(tx)
	users := NewUserRepository(tx)

	challenger := domain.NewUser(testID(1), "Challenger", "", "challenger")
	opponent := domain.NewUser(testID(2), "Opponent", "", "opponent")
	for _, u := range []*domain.User{challenger, opponent} {
		if err := users.Save(u); err != nil {
			t.Fatalf("save user: %v", err)
		}
	}
	duel := domain.NewDuel(testID(3), *challenger, *opponent)
	if err := repo.Store(duel); err != nil {
		t.Fatalf("store duel: %v", err)
	}

	duel.MarkAsDeclined()
	if claimed, err := repo.Resolve(duel); err != nil || !claimed {
		t.Fatalf("resolve: claimed=%v err=%v", claimed, err)
	}

	got, err := repo.GetByID(duel.ID)
	if err != nil || !got.IsDeclined() || got.IsPlayed() {
		t.Fatalf("expected a declined duel, got %+v err=%v", got, err)
	}
}

func TestDuelRepositoryHasPendingFrom(t *testing.T) {
	tx := withTx(t)
	repo := NewDuelRepository(tx)
	users := NewUserRepository(tx)

	challenger := domain.NewUser(testID(1), "Challenger", "", "challenger")
	opponent := domain.NewUser(testID(2), "Opponent", "", "opponent")
	for _, u := range []*domain.User{challenger, opponent} {
		if err := users.Save(u); err != nil {
			t.Fatalf("save user: %v", err)
		}
	}
	chatID := testID(3)
	since := time.Now().Add(-domain.DuelTTL)

	if has, err := repo.HasPendingFrom(chatID, challenger.ID, since); err != nil || has {
		t.Fatalf("expected nothing pending, got %v err=%v", has, err)
	}

	duel := domain.NewDuel(chatID, *challenger, *opponent)
	if err := repo.Store(duel); err != nil {
		t.Fatalf("store duel: %v", err)
	}
	if has, err := repo.HasPendingFrom(chatID, challenger.ID, since); err != nil || !has {
		t.Fatalf("expected a pending duel, got %v err=%v", has, err)
	}
	if has, _ := repo.HasPendingFrom(chatID, opponent.ID, since); has {
		t.Fatal("the opponent has no pending challenge of their own")
	}
	if has, _ := repo.HasPendingFrom(chatID, challenger.ID, time.Now().Add(time.Minute)); has {
		t.Fatal("a duel older than the window must not count")
	}

	duel.MarkAsDeclined()
	if _, err := repo.Resolve(duel); err != nil {
		t.Fatal(err)
	}
	if has, _ := repo.HasPendingFrom(chatID, challenger.ID, since); has {
		t.Fatal("a closed duel must not count")
	}
}
