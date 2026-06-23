package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
)

type fakeDuelRepo struct {
	domain.DuelRepository
	mu         sync.Mutex
	duel       *domain.Duel
	getErr     error
	resolveErr error
	resolved   int
}

func (r *fakeDuelRepo) GetByID(uuid.UUID) (*domain.Duel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.getErr != nil {
		return nil, r.getErr
	}
	cp := *r.duel
	return &cp, nil
}

func (r *fakeDuelRepo) Resolve(d *domain.Duel) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.resolveErr != nil {
		return false, r.resolveErr
	}
	if r.duel.IsClosed() {
		return false, nil
	}
	r.duel.PlayedAt, r.duel.DeclinedAt, r.duel.LoserID, r.duel.ProofURL = d.PlayedAt, d.DeclinedAt, d.LoserID, d.ProofURL
	r.resolved++
	return true, nil
}

type duelUC struct {
	uc   *RespondDuelUseCase
	repo *fakeDuelRepo
	duel *domain.Duel
	now  time.Time
}

func newDuelUC(t *testing.T, relayStatus int) *duelUC {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if relayStatus != http.StatusOK {
			w.WriteHeader(relayStatus)
			return
		}
		_ = json.NewEncoder(w).Encode(drand.Beacon{
			Round:      7,
			Randomness: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		})
	}))
	t.Cleanup(server.Close)

	duel := domain.NewDuel(100, *domain.NewUser(1, "Alice", "", "alice"), *domain.NewUser(2, "Bob", "", "bob"))
	repo := &fakeDuelRepo{duel: duel}

	return &duelUC{
		uc:   NewRespondDuelUseCase(service.NewLocker(), repo, drand.NewClientWithURL(server.URL)),
		repo: repo,
		duel: duel,
		now:  duel.CreatedAt.Add(time.Minute),
	}
}

func TestRespondDuel_AcceptStoresTheOutcomeAndReturnsAParticipantAsLoser(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)

	out, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now)
	if err != nil {
		t.Fatal(err)
	}

	if out.Loser.ID != 1 && out.Loser.ID != 2 {
		t.Fatalf("loser %d is not a participant", out.Loser.ID)
	}
	if !e.repo.duel.IsPlayed() || e.repo.duel.LoserID != out.Loser.ID || e.repo.duel.ProofURL == "" || out.ProofURL != e.repo.duel.ProofURL {
		t.Fatalf("outcome was not stored consistently: %+v", e.repo.duel)
	}
}

func TestRespondDuel_OutcomeIsDeterministicForTheSameBeaconAndDuel(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)

	first, _, err := e.uc.shoot(context.Background(), e.duel)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := e.uc.shoot(context.Background(), e.duel)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("same beacon and duel gave different losers: %d vs %d", first, second)
	}
}

func TestRespondDuel_BothParticipantsCanLose(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	losers := map[int64]bool{}

	for i := 0; i < 64 && len(losers) < 2; i++ {
		e.duel.ID = uuid.New()
		loser, _, err := e.uc.shoot(context.Background(), e.duel)
		if err != nil {
			t.Fatal(err)
		}
		losers[loser] = true
	}

	if len(losers) != 2 {
		t.Fatalf("expected both participants to be able to lose, got %v", losers)
	}
}

func TestRespondDuel_OnlyTheOpponentMayRespond(t *testing.T) {
	for _, uid := range []int64{1, 3} {
		e := newDuelUC(t, http.StatusOK)

		if _, err := e.uc.Accept(context.Background(), e.duel.ID, uid, e.now); !errors.Is(err, ErrDuelNotOpponent) {
			t.Fatalf("accept by %d: expected ErrDuelNotOpponent, got %v", uid, err)
		}
		if _, err := e.uc.Decline(e.duel.ID, uid, e.now); !errors.Is(err, ErrDuelNotOpponent) {
			t.Fatalf("decline by %d: expected ErrDuelNotOpponent, got %v", uid, err)
		}
		if e.repo.resolved != 0 {
			t.Fatal("nothing must be stored")
		}
	}
}

func TestRespondDuel_ExpiredDuelCannotBeAnswered(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	late := e.duel.CreatedAt.Add(domain.DuelTTL + time.Second)

	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, late); !errors.Is(err, ErrDuelExpired) {
		t.Fatalf("accept: expected ErrDuelExpired, got %v", err)
	}
	if _, err := e.uc.Decline(e.duel.ID, 2, late); !errors.Is(err, ErrDuelExpired) {
		t.Fatalf("decline: expected ErrDuelExpired, got %v", err)
	}
	if e.repo.resolved != 0 {
		t.Fatal("nothing must be stored")
	}
}

func TestRespondDuel_DeclineClosesTheDuel(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)

	if _, err := e.uc.Decline(e.duel.ID, 2, e.now); err != nil {
		t.Fatal(err)
	}
	if !e.repo.duel.IsDeclined() || e.repo.duel.IsPlayed() {
		t.Fatalf("expected a declined duel, got %+v", e.repo.duel)
	}
	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); !errors.Is(err, ErrDuelClosed) {
		t.Fatalf("a declined duel must not be playable, got %v", err)
	}
}

func TestRespondDuel_PlayedDuelCannotBePlayedAgain(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); err != nil {
		t.Fatal(err)
	}

	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); !errors.Is(err, ErrDuelClosed) {
		t.Fatalf("expected ErrDuelClosed, got %v", err)
	}
	if _, err := e.uc.Decline(e.duel.ID, 2, e.now); !errors.Is(err, ErrDuelClosed) {
		t.Fatalf("expected ErrDuelClosed, got %v", err)
	}
	if e.repo.resolved != 1 {
		t.Fatalf("expected exactly one stored outcome, got %d", e.repo.resolved)
	}
}

func TestRespondDuel_ConcurrentAcceptsStoreOneOutcome(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successes != 1 || e.repo.resolved != 1 {
		t.Fatalf("expected one winner of the race, got %d successes and %d stored outcomes", successes, e.repo.resolved)
	}
}

func TestRespondDuel_LosingTheClaimIsReportedAsClosed(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	closed := *e.duel
	closed.MarkAsPlayed(1, "x")
	e.repo.duel = &closed
	e.repo.getErr = nil
	stale := &staleRepo{fakeDuelRepo: e.repo, open: e.duel}
	uc := NewRespondDuelUseCase(service.NewLocker(), stale, e.uc.drand)

	if _, err := uc.Accept(context.Background(), e.duel.ID, 2, e.now); !errors.Is(err, ErrDuelClosed) {
		t.Fatalf("expected ErrDuelClosed, got %v", err)
	}
}

type staleRepo struct {
	*fakeDuelRepo
	open *domain.Duel
}

func (r *staleRepo) GetByID(uuid.UUID) (*domain.Duel, error) {
	cp := *r.open
	return &cp, nil
}

func TestRespondDuel_RandomnessFailureStoresNothing(t *testing.T) {
	e := newDuelUC(t, http.StatusInternalServerError)

	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); err == nil {
		t.Fatal("expected an error")
	}
	if e.repo.resolved != 0 || e.repo.duel.IsClosed() {
		t.Fatal("a failed shot must leave the duel open")
	}
}

func TestRespondDuel_StorageErrorsArePropagated(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	e.repo.resolveErr = errors.New("db down")

	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); !errors.Is(err, e.repo.resolveErr) {
		t.Fatalf("expected the storage error, got %v", err)
	}

	e.repo.getErr = errors.New("db down")
	if _, err := e.uc.Decline(e.duel.ID, 2, e.now); !errors.Is(err, e.repo.getErr) {
		t.Fatalf("expected the storage error, got %v", err)
	}
}

func TestRespondDuel_ConcurrentLockContentionIsReported(t *testing.T) {
	e := newDuelUC(t, http.StatusOK)
	locker := service.NewLocker()
	e.uc.locker = locker
	release, ok := locker.TryLock(duelLockKey(e.duel.ID))
	if !ok {
		t.Fatal("lock must be free")
	}
	defer release()

	if _, err := e.uc.Accept(context.Background(), e.duel.ID, 2, e.now); !errors.Is(err, service.ErrLockNotAcquired) {
		t.Fatalf("expected ErrLockNotAcquired, got %v", err)
	}
}
