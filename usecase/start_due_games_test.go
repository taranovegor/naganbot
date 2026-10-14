package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

type dueGamesRepo struct {
	fakeGameRepo

	games map[uuid.UUID]*domain.Game
	err   error
}

func (r *dueGamesRepo) GetByID(id uuid.UUID) (*domain.Game, error) { return r.games[id], nil }

func (r *dueGamesRepo) GetActiveCountdownGames(time.Time) ([]*domain.Game, error) {
	var games []*domain.Game
	for _, game := range r.games {
		games = append(games, game)
	}
	return games, r.err
}

// byGameGunslingerRepo gives each game its own gunslingers.
type byGameGunslingerRepo struct {
	fakeGunslingerRepo

	byGame map[uuid.UUID][]*domain.Gunslinger
}

func (r *byGameGunslingerRepo) GetByGameID(id uuid.UUID) ([]*domain.Gunslinger, error) {
	return r.byGame[id], nil
}

func countdownGame(chatID int64, players int, joinedAgo time.Duration) (*domain.Game, []*domain.Gunslinger) {
	game := &domain.Game{ID: uuid.Must(uuid.NewV7()), ChatID: chatID, Mode: domain.GameModeCountdown}
	for i := 0; i < players; i++ {
		gunslinger := &domain.Gunslinger{
			ID:       uuid.Must(uuid.NewV7()),
			GameID:   game.ID,
			Game:     game,
			JoinedAt: time.Now().Add(-joinedAgo),
		}
		game.Gunslingers = append(game.Gunslingers, gunslinger)
	}

	return game, game.Gunslingers
}

func newStartDueGames(t *testing.T, repo *dueGamesRepo, gunslingers *byGameGunslingerRepo) (*StartDueGamesUseCase, *fakeUnitOfWork) {
	t.Helper()

	uow := &fakeUnitOfWork{}
	play := NewPlayGameUseCase(service.NewLocker(), repo, gunslingers, uow, newTestNagan(t))

	return NewStartDueGamesUseCase(repo, play), uow
}

func TestStartDueGamesPlaysOnlyTheGamesWhoseTimeHasCome(t *testing.T) {
	due, dueGunslingers := countdownGame(1, 3, 2*time.Hour)
	waiting, waitingGunslingers := countdownGame(2, 3, time.Minute)
	tooSmall, tooSmallGunslingers := countdownGame(3, 2, 2*time.Hour)

	repo := &dueGamesRepo{games: map[uuid.UUID]*domain.Game{due.ID: due, waiting.ID: waiting, tooSmall.ID: tooSmall}}
	gunslingers := &byGameGunslingerRepo{byGame: map[uuid.UUID][]*domain.Gunslinger{
		due.ID:      dueGunslingers,
		waiting.ID:  waitingGunslingers,
		tooSmall.ID: tooSmallGunslingers,
	}}
	uc, uow := newStartDueGames(t, repo, gunslingers)

	started, err := uc.Execute(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(started) != 1 || started[0].ChatID != due.ChatID || started[0].Report == nil {
		t.Fatalf("expected only the due game to be started, got %+v", started)
	}
	if uow.committedFor != due {
		t.Fatal("expected the due game to be committed as played")
	}
}

func TestStartDueGamesIgnoresGamesSomeoneElseGotToFirst(t *testing.T) {
	played, gunslingers := countdownGame(1, 3, 2*time.Hour)
	played.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}

	repo := &dueGamesRepo{games: map[uuid.UUID]*domain.Game{played.ID: played}}
	uc, _ := newStartDueGames(t, repo, &byGameGunslingerRepo{byGame: map[uuid.UUID][]*domain.Gunslinger{played.ID: gunslingers}})

	started, err := uc.Execute(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("expected an already played game not to be an error, got %v", err)
	}
	if len(started) != 0 {
		t.Fatalf("expected nothing to be started, got %+v", started)
	}
}

func TestStartDueGamesKeepsGoingAfterAGameFails(t *testing.T) {
	failing, failingGunslingers := countdownGame(1, 3, 2*time.Hour)
	working, workingGunslingers := countdownGame(2, 3, 2*time.Hour)

	repo := &dueGamesRepo{games: map[uuid.UUID]*domain.Game{failing.ID: failing, working.ID: working}}
	gunslingers := &byGameGunslingerRepo{byGame: map[uuid.UUID][]*domain.Gunslinger{
		failing.ID: failingGunslingers,
		working.ID: workingGunslingers,
	}}

	commitErr := errors.New("commit failed")
	uow := &failingFor{failing: failing.ID, err: commitErr}
	play := NewPlayGameUseCase(service.NewLocker(), repo, gunslingers, uow, newTestNagan(t))
	uc := NewStartDueGamesUseCase(repo, play)

	started, err := uc.Execute(context.Background(), time.Now())
	if !errors.Is(err, commitErr) {
		t.Fatalf("expected the failure to be reported, got %v", err)
	}
	if len(started) != 1 || started[0].ChatID != working.ChatID {
		t.Fatalf("expected the healthy game to still start, got %+v", started)
	}
}

type failingFor struct {
	failing uuid.UUID
	err     error
}

func (u *failingFor) CommitPlayedGame(game *domain.Game, _ []*domain.Gunslinger) error {
	if game.ID == u.failing {
		return u.err
	}
	return nil
}

func TestStartDueGamesReturnsRepositoryErrors(t *testing.T) {
	boom := errors.New("db is down")
	uc, _ := newStartDueGames(t, &dueGamesRepo{err: boom}, &byGameGunslingerRepo{})

	if _, err := uc.Execute(context.Background(), time.Now()); !errors.Is(err, boom) {
		t.Fatalf("expected the repository error, got %v", err)
	}
}
