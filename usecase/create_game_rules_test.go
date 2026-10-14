package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"gorm.io/gorm"
)

type scriptedGameRepo struct {
	lifecycleGameRepo

	active      *domain.Game
	activeErr   error
	playedToday bool

	idleSince time.Time
	since     time.Time
	stored    *domain.Game
}

func (r *scriptedGameRepo) GetActiveForChat(_ int64, idleSince time.Time) (*domain.Game, error) {
	r.idleSince = idleSince
	if r.activeErr != nil {
		return nil, r.activeErr
	}
	if r.active == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return r.active, nil
}

func (r *scriptedGameRepo) HasPlayedSince(_ int64, since time.Time) bool {
	r.since = since
	return r.playedToday
}

func (r *scriptedGameRepo) Store(game *domain.Game) error {
	r.stored = game
	return nil
}

func newScriptedCreateGame(repo *scriptedGameRepo) *CreateGameUseCase {
	chat := domain.Chat{ID: 100}
	chat.Settings.RequiredPlayers = 6

	return NewCreateGameUseCase(
		service.NewLocker(),
		&lifecycleChatRepo{chat: chat},
		repo,
		&lifecycleUserRepo{users: map[int64]domain.User{7: {ID: 7}}},
	)
}

func TestCreateGameReturnsTheActiveGameWithoutCreatingAnother(t *testing.T) {
	active := &domain.Game{ID: uuid.Must(uuid.NewV7())}
	repo := &scriptedGameRepo{active: active, playedToday: true}

	game, err := newScriptedCreateGame(repo).Execute(100, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if game != active {
		t.Fatal("expected the active game to be returned")
	}
	if repo.stored != nil {
		t.Fatal("expected no new game to be stored")
	}
}

func TestCreateGameIgnoresLobbiesIdleForLongerThanTheTTL(t *testing.T) {
	repo := &scriptedGameRepo{}

	before := time.Now()
	if _, err := newScriptedCreateGame(repo).Execute(100, 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := before.Add(-domain.LobbyIdleTTL)
	if repo.idleSince.Before(expected) || repo.idleSince.After(expected.Add(time.Minute)) {
		t.Fatalf("expected the idle cutoff to be about %s, got %s", expected, repo.idleSince)
	}
	if repo.stored == nil {
		t.Fatal("expected a new game to be stored when the old lobby is idle")
	}
}

func TestCreateGameRefusesWhileTheChatCoolsDownAfterTodaysGame(t *testing.T) {
	repo := &scriptedGameRepo{playedToday: true}

	_, err := newScriptedCreateGame(repo).Execute(100, 7)
	if !errors.Is(err, ErrGameCooldown) {
		t.Fatalf("expected ErrGameCooldown, got %v", err)
	}
	if repo.stored != nil {
		t.Fatal("expected no game to be stored during the cooldown")
	}

	now := time.Now()
	y, m, d := now.Date()
	if midnight := time.Date(y, m, d, 0, 0, 0, 0, now.Location()); !repo.since.Equal(midnight) && !repo.since.Equal(midnight.AddDate(0, 0, 1)) {
		t.Fatalf("expected the cooldown to be measured from today's midnight, got %s", repo.since)
	}
}

func TestCreateGamePropagatesRepositoryErrors(t *testing.T) {
	boom := errors.New("db is down")
	repo := &scriptedGameRepo{activeErr: boom}

	_, err := newScriptedCreateGame(repo).Execute(100, 7)
	if !errors.Is(err, boom) {
		t.Fatalf("expected the repository error, got %v", err)
	}
}
