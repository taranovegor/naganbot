package app

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/usecase"
)

type fakeStarter struct {
	mu      sync.Mutex
	calls   int
	started []usecase.StartedGame
	err     error
}

func (s *fakeStarter) Execute(context.Context, time.Time) ([]usecase.StartedGame, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return s.started, s.err
}

func (s *fakeStarter) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

type fakeAnnouncer struct {
	mu     sync.Mutex
	chatID []int64
}

func (a *fakeAnnouncer) AnnounceResult(chatID int64, _ *service.HitReport) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.chatID = append(a.chatID, chatID)
}

func TestStartDueGamesAnnouncesEveryStartedGame(t *testing.T) {
	starter := &fakeStarter{started: []usecase.StartedGame{
		{ChatID: 1, Report: &service.HitReport{}},
		{ChatID: 2, Report: &service.HitReport{}},
	}}
	announcer := &fakeAnnouncer{}

	startDueGames(context.Background(), time.Now(), starter, announcer)

	if len(announcer.chatID) != 2 || announcer.chatID[0] != 1 || announcer.chatID[1] != 2 {
		t.Fatalf("expected both chats to be announced, got %v", announcer.chatID)
	}
}

func TestStartDueGamesStillAnnouncesWhatStartedDespiteAnError(t *testing.T) {
	starter := &fakeStarter{
		started: []usecase.StartedGame{{ChatID: 1, Report: &service.HitReport{}}},
		err:     errors.New("one game failed"),
	}
	announcer := &fakeAnnouncer{}

	startDueGames(context.Background(), time.Now(), starter, announcer)

	if len(announcer.chatID) != 1 {
		t.Fatalf("expected the started game to be announced, got %v", announcer.chatID)
	}
}

func TestRunSchedulerChecksImmediatelyThenOnEveryTickUntilCancelled(t *testing.T) {
	starter := &fakeStarter{}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		runScheduler(ctx, 5*time.Millisecond, starter, &fakeAnnouncer{})
	}()

	deadline := time.After(2 * time.Second)
	for starter.callCount() < 3 {
		select {
		case <-deadline:
			t.Fatalf("expected repeated checks, got %d", starter.callCount())
		case <-time.After(time.Millisecond):
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected the scheduler to stop once the context is cancelled")
	}
}
