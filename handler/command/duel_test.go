package command

import (
	"context"
	"errors"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type duelCmdBot struct {
	messages  []string
	keyboards []service.Keyboard
	kicked    []int64
	kickErr   error
	isBot     bool
	isMember  bool
}

func (b *duelCmdBot) SendMessage(_ int64, text string) { b.messages = append(b.messages, text) }
func (b *duelCmdBot) SendInlineKeyboard(_ int64, text string, kb service.Keyboard) {
	b.messages = append(b.messages, text)
	b.keyboards = append(b.keyboards, kb)
}
func (b *duelCmdBot) Kick(_ int64, userID int64) error {
	b.kicked = append(b.kicked, userID)
	return b.kickErr
}
func (b *duelCmdBot) DeleteMessage(int64, int)       {}
func (b *duelCmdBot) IsBot(int64, int64) bool        { return b.isBot }
func (b *duelCmdBot) IsChatMember(int64, int64) bool { return b.isMember }

type duelCmdUsers struct {
	domain.UserRepository
	byID       map[int64]domain.User
	byUsername map[string]domain.User
}

func (r duelCmdUsers) Get(id int64) (domain.User, error) {
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return domain.User{}, errors.New("not found")
}
func (r duelCmdUsers) GetByUsername(name string) (domain.User, error) {
	if u, ok := r.byUsername[name]; ok {
		return u, nil
	}
	return domain.User{}, errors.New("not found")
}

type duelCmdRepo struct {
	domain.DuelRepository
	stored     []*domain.Duel
	storeErr   error
	pending    bool
	pendingErr error
}

func (r *duelCmdRepo) HasPendingFrom(int64, int64, time.Time) (bool, error) {
	return r.pending, r.pendingErr
}

func (r *duelCmdRepo) Store(d *domain.Duel) error {
	if r.storeErr != nil {
		return r.storeErr
	}
	r.stored = append(r.stored, d)
	return nil
}

type duelCmdEnv struct {
	bot   *duelCmdBot
	repo  *duelCmdRepo
	trans *translator.Translator
	hdlr  Handler
}

func newDuelCmdEnv() *duelCmdEnv {
	alice := domain.NewUser(1, "Alice", "", "alice")
	bob := domain.NewUser(2, "Bob", "", "bob")
	bot := &duelCmdBot{isMember: true}
	repo := &duelCmdRepo{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	users := duelCmdUsers{
		byID:       map[int64]domain.User{1: *alice, 2: *bob},
		byUsername: map[string]domain.User{"bob": *bob},
	}
	return &duelCmdEnv{bot: bot, repo: repo, trans: trans, hdlr: NewDuelHandler(bot, trans, users, repo)}
}

func (e *duelCmdEnv) text(key string) string { return e.trans.Get(key, translator.Config{}) }

func duelMessage(from int64, args string, reply *tgbotapi.User) *tgbotapi.Message {
	text := "/ngduel"
	if args != "" {
		text += " " + args
	}
	msg := &tgbotapi.Message{
		Text:     text,
		Chat:     &tgbotapi.Chat{ID: 100},
		From:     &tgbotapi.User{ID: from},
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len("/ngduel")}},
	}
	if reply != nil {
		msg.ReplyToMessage = &tgbotapi.Message{From: reply}
	}
	return msg
}

func (e *duelCmdEnv) expectSingleMessage(t *testing.T, want string) {
	t.Helper()
	if len(e.bot.messages) != 1 || e.bot.messages[0] != want {
		t.Fatalf("expected single message %q, got %q", want, e.bot.messages)
	}
	if len(e.repo.stored) != 0 {
		t.Fatal("no duel must be stored")
	}
}

func TestDuelHandler_Name(t *testing.T) {
	if got := newDuelCmdEnv().hdlr.Name(); got != "duel" {
		t.Fatalf("unexpected name %q", got)
	}
}

func TestDuelHandler_ChallengeByUsernameStoresDuelAndSendsKeyboard(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	if len(e.repo.stored) != 1 {
		t.Fatalf("expected one stored duel, got %d", len(e.repo.stored))
	}
	duel := e.repo.stored[0]
	if duel.ChallengerID != 1 || duel.OpponentID != 2 || duel.ChatID != 100 {
		t.Fatalf("unexpected duel: %+v", duel)
	}
	if len(e.bot.keyboards) != 1 || len(e.bot.keyboards[0]) != 1 || len(e.bot.keyboards[0][0]) != 2 {
		t.Fatalf("expected one row with two buttons, got %+v", e.bot.keyboards)
	}
	row := e.bot.keyboards[0][0]
	if want := callback.DuelResponse.SetArgs(duel.ID.String(), "accept").ToString(); row[0].Data != want {
		t.Fatalf("accept data: want %q, got %q", want, row[0].Data)
	}
	if want := callback.DuelResponse.SetArgs(duel.ID.String(), "decline").ToString(); row[1].Data != want {
		t.Fatalf("decline data: want %q, got %q", want, row[1].Data)
	}
	if row[0].Text != e.text("duel accept button") || row[1].Text != e.text("duel decline button") {
		t.Fatalf("unexpected button labels: %+v", row)
	}
	if _, err := uuid.Parse(duel.ID.String()); err != nil {
		t.Fatal(err)
	}
}

func TestDuelHandler_ChallengeByReplyDoesNotNeedStoredUser(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(1, "", &tgbotapi.User{ID: 77, FirstName: "Carl", UserName: "carl"}))

	if len(e.repo.stored) != 1 || e.repo.stored[0].OpponentID != 77 || e.repo.stored[0].Opponent.Username != "carl" {
		t.Fatalf("unexpected stored duels: %+v", e.repo.stored)
	}
}

func TestDuelHandler_ReplyToBotIsRejected(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(1, "", &tgbotapi.User{ID: 77, IsBot: true}))

	e.expectSingleMessage(t, e.text("duel is bot"))
}

func TestDuelHandler_NoTarget(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(1, "", nil))

	e.expectSingleMessage(t, e.text("duel no target"))
}

func TestDuelHandler_UnknownUsername(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(1, "@ghost", nil))

	e.expectSingleMessage(t, e.text("duel unknown user"))
}

func TestDuelHandler_SelfDuelKicksChallenger(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(2, "@bob", nil))

	if len(e.bot.kicked) != 1 || e.bot.kicked[0] != 2 {
		t.Fatalf("expected challenger kicked, got %v", e.bot.kicked)
	}
	if len(e.bot.messages) != 0 || len(e.repo.stored) != 0 {
		t.Fatal("nothing else must happen")
	}
}

func TestDuelHandler_SelfDuelFallsBackToMessageWhenKickFails(t *testing.T) {
	e := newDuelCmdEnv()
	e.bot.kickErr = errors.New("not enough rights")

	e.hdlr.Execute(context.Background(), duelMessage(2, "@bob", nil))

	e.expectSingleMessage(t, e.text("duel self"))
}

func TestDuelHandler_OpponentNotInChat(t *testing.T) {
	e := newDuelCmdEnv()
	e.bot.isMember = false

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	e.expectSingleMessage(t, e.text("duel not in chat"))
}

func TestDuelHandler_OpponentIsBot(t *testing.T) {
	e := newDuelCmdEnv()
	e.bot.isBot = true

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	e.expectSingleMessage(t, e.text("duel is bot"))
}

func TestDuelHandler_UnknownChallenger(t *testing.T) {
	e := newDuelCmdEnv()

	e.hdlr.Execute(context.Background(), duelMessage(999, "@bob", nil))

	e.expectSingleMessage(t, e.text("something went wrong"))
}

func TestDuelHandler_StoreFailure(t *testing.T) {
	e := newDuelCmdEnv()
	e.repo.storeErr = errors.New("db down")

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	e.expectSingleMessage(t, e.text("something went wrong"))
	if len(e.bot.keyboards) != 0 {
		t.Fatal("no keyboard must be sent when the duel was not stored")
	}
}

func TestDuelHandler_RejectsASecondChallengeWhileOneIsPending(t *testing.T) {
	e := newDuelCmdEnv()
	e.repo.pending = true

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	e.expectSingleMessage(t, e.text("duel pending"))
}

func TestDuelHandler_PendingLookupFailure(t *testing.T) {
	e := newDuelCmdEnv()
	e.repo.pendingErr = errors.New("db down")

	e.hdlr.Execute(context.Background(), duelMessage(1, "@bob", nil))

	e.expectSingleMessage(t, e.text("something went wrong"))
}
