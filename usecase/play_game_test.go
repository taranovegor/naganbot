package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
)

func TestGameLockKeyIsUniquePerGame(t *testing.T) {
	first := uuid.Must(uuid.NewV7())
	second := uuid.Must(uuid.NewV7())

	locker := service.NewLocker()
	firstMutex := locker.LockFor(gameLockKey(first))
	secondMutex := locker.LockFor(gameLockKey(second))

	if firstMutex == secondMutex {
		t.Fatal("two distinct games were assigned the same lock")
	}
}

type fakeGameRepo struct {
	game *domain.Game
}

func (r *fakeGameRepo) GetByID(uuid.UUID) (*domain.Game, error)                { return r.game, nil }
func (r *fakeGameRepo) GetLatestGamesInChat(int64, int) ([]domain.Game, error) { return nil, nil }
func (r *fakeGameRepo) GetLatestForChat(int64) (*domain.Game, error)           { return nil, nil }
func (r *fakeGameRepo) GetActiveForChat(int64) (*domain.Game, error)           { return nil, nil }
func (r *fakeGameRepo) Store(*domain.Game) error                               { return nil }
func (r *fakeGameRepo) HasActiveOrCreatedTodayInChat(int64) bool               { return false }

type fakeGunslingerRepo struct {
	gunslingers []*domain.Gunslinger
}

func (r *fakeGunslingerRepo) Store(*domain.Gunslinger) error { return nil }
func (r *fakeGunslingerRepo) GetByGameID(uuid.UUID) ([]*domain.Gunslinger, error) {
	return r.gunslingers, nil
}
func (r *fakeGunslingerRepo) IsPlayerExistsInGame(int64, uuid.UUID) bool { return false }
func (r *fakeGunslingerRepo) GetTopShotPlayersInChat(int64) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r *fakeGunslingerRepo) GetTopShotPlayersByYearInChat(int64, int) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r *fakeGunslingerRepo) CountNumberOfPlayerGamesInChat(int64, int64) int64 { return 0 }
func (r *fakeGunslingerRepo) CountNumberOfSelfShotsInChat(int64, int64) int64   { return 0 }

type fakeUnitOfWork struct {
	err          error
	committedFor *domain.Game
}

func (u *fakeUnitOfWork) CommitPlayedGame(game *domain.Game, _ []*domain.Gunslinger) error {
	u.committedFor = game
	return u.err
}

func newTestNagan(t *testing.T) *service.Nagan {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(drand.Beacon{
			Round:      1,
			Randomness: "aa",
		})
	}))
	t.Cleanup(server.Close)

	bulletFactory := service.NewBulletFactory(service.NewLeadBullet())
	return service.NewNagan(bulletFactory, drand.NewClientWithURL(server.URL))
}

func TestExecuteDoesNotPersistAnythingWhenTheAtomicCommitFails(t *testing.T) {
	gameID := uuid.Must(uuid.NewV7())
	game := &domain.Game{ID: gameID, PlayersCount: 2}
	gunslingers := []*domain.Gunslinger{
		{ID: uuid.Must(uuid.NewV7()), GameID: gameID, Game: game},
		{ID: uuid.Must(uuid.NewV7()), GameID: gameID, Game: game},
	}

	gameRepo := &fakeGameRepo{game: game}
	gunslingerRepo := &fakeGunslingerRepo{gunslingers: gunslingers}
	commitErr := errors.New("gunslinger update failed")
	uow := &fakeUnitOfWork{err: commitErr}

	uc := NewPlayGameUseCase(service.NewLocker(), gameRepo, gunslingerRepo, uow, newTestNagan(t))

	_, err := uc.Execute(context.Background(), gameID)

	if !errors.Is(err, commitErr) {
		t.Fatalf("expected the commit error to be returned, got %v", err)
	}
	if uow.committedFor == nil || !uow.committedFor.IsPlayed() {
		t.Fatal("expected the unit of work to receive the game marked as played")
	}
}
