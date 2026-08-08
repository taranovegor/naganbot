package usecase

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
)

type joinGameGunslingerRepo struct {
	exists   bool
	storeErr error
	stored   *domain.Gunslinger
}

func (r *joinGameGunslingerRepo) Store(g *domain.Gunslinger) error {
	r.stored = g
	return r.storeErr
}
func (r *joinGameGunslingerRepo) GetByGameID(uuid.UUID) ([]*domain.Gunslinger, error) {
	return nil, nil
}
func (r *joinGameGunslingerRepo) IsPlayerExistsInGame(int64, uuid.UUID) bool { return r.exists }
func (r *joinGameGunslingerRepo) GetTopShotPlayersInChat(int64) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r *joinGameGunslingerRepo) GetTopShotPlayersByYearInChat(int64, int) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r *joinGameGunslingerRepo) CountNumberOfPlayerGamesInChat(int64, int64) int64 { return 0 }
func (r *joinGameGunslingerRepo) CountNumberOfSelfShotsInChat(int64, int64) int64   { return 0 }

type joinGameUserRepo struct {
	user   domain.User
	getErr error
}

func (r *joinGameUserRepo) Get(int64) (domain.User, error)          { return r.user, r.getErr }
func (r *joinGameUserRepo) GetByIDs([]int64) ([]domain.User, error) { return nil, nil }
func (r *joinGameUserRepo) Save(*domain.User) error                 { return nil }

func TestJoinGameUseCaseReturnsErrWhenPlayerAlreadyInGame(t *testing.T) {
	gunslingerRepo := &joinGameGunslingerRepo{exists: true}
	userRepo := &joinGameUserRepo{}

	uc := NewJoinGameUseCase(nil, gunslingerRepo, userRepo)

	_, err := uc.Execute(uuid.Must(uuid.NewV7()), 7)

	if !errors.Is(err, ErrPlayerAlreadyInGame) {
		t.Fatalf("expected ErrPlayerAlreadyInGame, got %v", err)
	}
	if gunslingerRepo.stored != nil {
		t.Fatal("expected no gunslinger to be stored")
	}
}

func TestJoinGameUseCasePropagatesUserLookupError(t *testing.T) {
	getErr := errors.New("connection refused")
	gunslingerRepo := &joinGameGunslingerRepo{}
	userRepo := &joinGameUserRepo{getErr: getErr}

	uc := NewJoinGameUseCase(nil, gunslingerRepo, userRepo)

	_, err := uc.Execute(uuid.Must(uuid.NewV7()), 7)

	if !errors.Is(err, getErr) {
		t.Fatalf("expected the user lookup error to be returned, got %v", err)
	}
	if gunslingerRepo.stored != nil {
		t.Fatal("expected no gunslinger to be stored")
	}
}

func TestJoinGameUseCasePropagatesStoreError(t *testing.T) {
	storeErr := errors.New("write failed")
	gunslingerRepo := &joinGameGunslingerRepo{storeErr: storeErr}
	userRepo := &joinGameUserRepo{user: domain.User{ID: 7, FirstName: "Kate"}}

	uc := NewJoinGameUseCase(nil, gunslingerRepo, userRepo)

	_, err := uc.Execute(uuid.Must(uuid.NewV7()), 7)

	if !errors.Is(err, storeErr) {
		t.Fatalf("expected the store error to be returned, got %v", err)
	}
}

func TestJoinGameUseCaseStoresAndReturnsTheGunslingerOnSuccess(t *testing.T) {
	gameID := uuid.Must(uuid.NewV7())
	user := domain.User{ID: 7, FirstName: "Kate"}
	gunslingerRepo := &joinGameGunslingerRepo{}
	userRepo := &joinGameUserRepo{user: user}

	uc := NewJoinGameUseCase(nil, gunslingerRepo, userRepo)

	gunslinger, err := uc.Execute(gameID, 7)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gunslinger == nil {
		t.Fatal("expected a gunslinger to be returned")
	}
	if gunslinger.GameID != gameID {
		t.Fatalf("expected GameID %s, got %s", gameID, gunslinger.GameID)
	}
	if gunslinger.PlayerID != 7 || gunslinger.Player != user {
		t.Fatalf("expected the gunslinger to belong to %+v, got PlayerID=%d Player=%+v", user, gunslinger.PlayerID, gunslinger.Player)
	}
	if gunslingerRepo.stored != gunslinger {
		t.Fatal("expected the gunslinger to be stored")
	}
}
