package callback

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type duelBot struct {
	answers []string
	edits   []string
	replies []string
	kicked  []int64
	kickErr error
}

func (b *duelBot) AnswerCallback(_ string, text string)                { b.answers = append(b.answers, text) }
func (b *duelBot) EditMessageReplyMarkup(int64, int, service.Keyboard) {}
func (b *duelBot) IsAdmin(int64, int64) (bool, error)                  { return false, nil }
func (b *duelBot) EditMessageText(_ int64, _ int, text string) {
	b.edits = append(b.edits, text)
}
func (b *duelBot) ReplyMessage(_ int64, _ int, text string) { b.replies = append(b.replies, text) }
func (b *duelBot) Kick(_ int64, userID int64) error {
	b.kicked = append(b.kicked, userID)
	return b.kickErr
}

type duelRepoFake struct {
	domain.DuelRepository
	duel      *domain.Duel
	getErr    error
	resolveOK bool
	resolved  int
}

func (r *duelRepoFake) GetByID(uuid.UUID) (*domain.Duel, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	cp := *r.duel
	return &cp, nil
}

func (r *duelRepoFake) Resolve(d *domain.Duel) (bool, error) {
	if r.duel.IsClosed() {
		return false, nil
	}
	*r.duel = *d
	r.resolved++
	return true, nil
}

type duelEnv struct {
	bot   *duelBot
	repo  *duelRepoFake
	trans *translator.Translator
	hdlr  Handler
	duel  *domain.Duel
}

func newDuelEnv(t *testing.T, relayStatus int) *duelEnv {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	alice := domain.NewUser(1, "Alice", "", "alice")
	bob := domain.NewUser(2, "Bob", "", "bob")
	duel := domain.NewDuel(100, *alice, *bob)

	bot := &duelBot{}
	repo := &duelRepoFake{duel: duel}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	uc := usecase.NewRespondDuelUseCase(service.NewLocker(), repo, drand.NewClientWithURL(server.URL))

	return &duelEnv{bot: bot, repo: repo, trans: trans, hdlr: NewDuelResponse(bot, trans, uc), duel: duel}
}

func (e *duelEnv) query(from int64, action string) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		ID:      "cb",
		Data:    DuelResponse.SetArgs(e.duel.ID.String(), action).ToString(),
		From:    &tgbotapi.User{ID: from},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 100}, MessageID: 5},
	}
}

func (e *duelEnv) text(key string) string { return e.trans.Get(key, translator.Config{}) }

func (e *duelEnv) expectOnlyAnswer(t *testing.T, want string) {
	t.Helper()
	if len(e.bot.answers) != 1 || e.bot.answers[0] != want {
		t.Fatalf("expected single answer %q, got %q", want, e.bot.answers)
	}
	if len(e.bot.edits)+len(e.bot.replies)+len(e.bot.kicked)+e.repo.resolved != 0 {
		t.Fatalf("no side effects expected: %+v resolved=%d", e.bot, e.repo.resolved)
	}
}

func TestDuelResponse_Pattern(t *testing.T) {
	if got := newDuelEnv(t, http.StatusOK).hdlr.Pattern(); got != DuelResponse {
		t.Fatalf("unexpected pattern %q", got)
	}
}

func TestDuelResponse_InvalidDuelID(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)
	q := e.query(2, "accept")
	q.Data = DuelResponse.SetArgs("not-a-uuid", "accept").ToString()

	e.hdlr.Execute(context.Background(), q)

	e.expectOnlyAnswer(t, e.text("something went wrong"))
}

func TestDuelResponse_DuelNotFound(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)
	e.repo.getErr = errors.New("not found")

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	e.expectOnlyAnswer(t, e.text("something went wrong"))
}

func TestDuelResponse_AlreadyPlayedIsIgnored(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)
	e.duel.MarkAsPlayed(1, "proof")

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	e.expectOnlyAnswer(t, "")
}

func TestDuelResponse_OnlyOpponentMayRespond(t *testing.T) {
	for _, action := range []string{"accept", "decline"} {
		for _, from := range []int64{1, 3} {
			e := newDuelEnv(t, http.StatusOK)

			e.hdlr.Execute(context.Background(), e.query(from, action))

			e.expectOnlyAnswer(t, e.text("duel not your"))
		}
	}
}

func TestDuelResponse_UnknownActionDoesNothing(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)

	e.hdlr.Execute(context.Background(), e.query(2, "bogus"))

	if len(e.bot.answers)+len(e.bot.edits)+len(e.bot.replies)+len(e.bot.kicked)+e.repo.resolved != 0 {
		t.Fatalf("unexpected side effects: %+v", e.bot)
	}
}

func TestDuelResponse_DeclineEditsMessageAndClosesTheDuel(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)

	e.hdlr.Execute(context.Background(), e.query(2, "decline"))

	if len(e.bot.edits) != 1 {
		t.Fatalf("expected one edit, got %q", e.bot.edits)
	}
	if len(e.bot.kicked) != 0 || !e.duel.IsDeclined() || e.duel.IsPlayed() {
		t.Fatal("declined duel must not kick anyone and must be stored as declined")
	}
}

func TestDuelResponse_AcceptPlaysDuelAndKicksTheLoser(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	if len(e.bot.edits) != 1 {
		t.Fatalf("expected the challenge message to be edited once, got %q", e.bot.edits)
	}
	if !e.duel.IsPlayed() || e.repo.resolved != 1 {
		t.Fatal("expected the duel to be marked as played and stored")
	}
	loser := e.duel.LoserID
	if loser != e.duel.ChallengerID && loser != e.duel.OpponentID {
		t.Fatalf("loser %d is not a participant", loser)
	}
	if len(e.bot.kicked) != 1 || e.bot.kicked[0] != loser {
		t.Fatalf("expected loser %d kicked, got %v", loser, e.bot.kicked)
	}
	if len(e.bot.replies) != 1 {
		t.Fatalf("expected one result reply, got %q", e.bot.replies)
	}
	if e.duel.ProofURL == "" {
		t.Fatal("expected a proof url")
	}
}

func TestDuelResponse_RandomnessFailureReportsErrorAndDoesNotKick(t *testing.T) {
	e := newDuelEnv(t, http.StatusInternalServerError)

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	if len(e.bot.answers) != 1 || e.bot.answers[0] != e.text("something went wrong") {
		t.Fatalf("expected error answer, got %q", e.bot.answers)
	}
	if len(e.bot.kicked) != 0 || e.repo.resolved != 0 || e.duel.IsPlayed() {
		t.Fatal("a failed shot must not kick anyone or mark the duel played")
	}
}

func TestDuelResponse_ExpiredDuelEditsTheMessage(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)
	e.duel.CreatedAt = time.Now().Add(-domain.DuelTTL - time.Minute)

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	if len(e.bot.edits) != 1 || e.bot.edits[0] != e.text("duel expired") {
		t.Fatalf("expected the expiry text, got %q", e.bot.edits)
	}
	if len(e.bot.kicked) != 0 || e.repo.resolved != 0 {
		t.Fatal("an expired duel must not be played")
	}
}

func TestDuelResponse_SecondAcceptIsIgnored(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))
	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	if e.repo.resolved != 1 || len(e.bot.kicked) != 1 || len(e.bot.replies) != 1 {
		t.Fatalf("expected a single outcome, got resolved=%d kicked=%v replies=%q", e.repo.resolved, e.bot.kicked, e.bot.replies)
	}
}

func TestDuelResponse_FailedKickIsAnnounced(t *testing.T) {
	e := newDuelEnv(t, http.StatusOK)
	e.bot.kickErr = errors.New("not enough rights")

	e.hdlr.Execute(context.Background(), e.query(2, "accept"))

	if len(e.bot.replies) != 2 || e.bot.replies[1] != e.text("duel kick failed") {
		t.Fatalf("expected the kick failure notice, got %q", e.bot.replies)
	}
}
