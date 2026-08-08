package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
	"gorm.io/gorm"
)

type joinTestWorld struct {
	mu sync.Mutex

	chat  domain.Chat
	users map[int64]domain.User

	game        *domain.Game
	gunslingers []*domain.Gunslinger

	commitErr error
}

type joinTestChatRepo struct{ w *joinTestWorld }

func (r joinTestChatRepo) Get(int64) (domain.Chat, error) {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	return r.w.chat, nil
}
func (r joinTestChatRepo) Save(*domain.Chat) error           { return nil }
func (r joinTestChatRepo) UpdateSettings(*domain.Chat) error { return nil }

type joinTestUserRepo struct{ w *joinTestWorld }

func (r joinTestUserRepo) Get(id int64) (domain.User, error) {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	return r.w.users[id], nil
}
func (r joinTestUserRepo) GetByIDs([]int64) ([]domain.User, error) { return nil, nil }
func (r joinTestUserRepo) Save(*domain.User) error                 { return nil }

type joinTestGameRepo struct{ w *joinTestWorld }

func (r joinTestGameRepo) GetByID(uuid.UUID) (*domain.Game, error) {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	return r.w.game, nil
}
func (r joinTestGameRepo) GetLatestGamesInChat(int64, int) ([]domain.Game, error) { return nil, nil }
func (r joinTestGameRepo) GetLatestForChat(int64) (*domain.Game, error)           { return nil, nil }
func (r joinTestGameRepo) GetActiveForChat(int64) (*domain.Game, error) {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	if r.w.game == nil || r.w.game.IsPlayed() {
		return nil, gorm.ErrRecordNotFound
	}
	return r.w.game, nil
}
func (r joinTestGameRepo) Store(game *domain.Game) error {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	r.w.game = game
	return nil
}
func (r joinTestGameRepo) HasActiveOrCreatedTodayInChat(int64) bool {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	return r.w.game != nil
}

type joinTestGunslingerRepo struct{ w *joinTestWorld }

func (r joinTestGunslingerRepo) Store(g *domain.Gunslinger) error {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	r.w.gunslingers = append(r.w.gunslingers, g)
	return nil
}
func (r joinTestGunslingerRepo) GetByGameID(uuid.UUID) ([]*domain.Gunslinger, error) {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	for _, g := range r.w.gunslingers {
		g.Game = r.w.game
	}
	return r.w.gunslingers, nil
}
func (r joinTestGunslingerRepo) IsPlayerExistsInGame(userID int64, _ uuid.UUID) bool {
	r.w.mu.Lock()
	defer r.w.mu.Unlock()
	for _, g := range r.w.gunslingers {
		if g.PlayerID == userID {
			return true
		}
	}
	return false
}
func (r joinTestGunslingerRepo) GetTopShotPlayersInChat(int64) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r joinTestGunslingerRepo) GetTopShotPlayersByYearInChat(int64, int) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}
func (r joinTestGunslingerRepo) CountNumberOfPlayerGamesInChat(int64, int64) int64 { return 0 }
func (r joinTestGunslingerRepo) CountNumberOfSelfShotsInChat(int64, int64) int64   { return 0 }

type joinTestUnitOfWork struct{ w *joinTestWorld }

func (u joinTestUnitOfWork) CommitPlayedGame(game *domain.Game, _ []*domain.Gunslinger) error {
	if u.w.commitErr != nil {
		return u.w.commitErr
	}
	u.w.mu.Lock()
	defer u.w.mu.Unlock()
	u.w.game = game
	return nil
}

type joinTestBot struct {
	mu             sync.Mutex
	messages       []string
	kicked         []int64
	kickErr        error
	deletedMessage int
	deleteCalled   bool
}

func (b *joinTestBot) SendMessage(_ int64, text string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.messages = append(b.messages, text)
}
func (b *joinTestBot) SendInlineKeyboard(int64, string, service.Keyboard) {}
func (b *joinTestBot) Kick(_ int64, userID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.kicked = append(b.kicked, userID)
	return b.kickErr
}
func (b *joinTestBot) DeleteMessage(_ int64, messageID int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.deleteCalled = true
	b.deletedMessage = messageID
}

func (b *joinTestBot) sentMessages() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.messages...)
}

func (b *joinTestBot) kickedPlayers() []int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]int64(nil), b.kicked...)
}

func possibleTexts(trans *translator.Translator, key string, cfg translator.Config) map[string]bool {
	texts := make(map[string]bool)
	for i := 1; ; i++ {
		cfg.OneOfMany = i
		got := trans.Get(key, cfg)
		if got == key {
			break
		}
		texts[got] = true
	}
	return texts
}

func newTestNaganWithBullet(t *testing.T, bullet service.Bullet) *service.Nagan {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(drand.Beacon{
			Round:      1,
			Randomness: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		})
	}))
	t.Cleanup(server.Close)

	factory := service.NewBulletFactory(bullet)
	return service.NewNagan(factory, drand.NewClientWithURL(server.URL))
}

func newJoinTestHandler(w *joinTestWorld, bot *joinTestBot, nagan *service.Nagan) *joinHandler {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	chatRepo := joinTestChatRepo{w}
	userRepo := joinTestUserRepo{w}
	gameRepo := joinTestGameRepo{w}
	gunslingerRepo := joinTestGunslingerRepo{w}
	uow := joinTestUnitOfWork{w}

	return &joinHandler{
		bot:          bot,
		createGameUC: usecase.NewCreateGameUseCase(service.NewLocker(), chatRepo, gameRepo, userRepo),
		joinGameUC:   usecase.NewJoinGameUseCase(gameRepo, gunslingerRepo, userRepo),
		playGameUC:   usecase.NewPlayGameUseCase(service.NewLocker(), gameRepo, gunslingerRepo, uow, nagan),
		trans:        trans,
		messageDelay: 0,
	}
}

func joinMessage(chatID int64, fromID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		From: &tgbotapi.User{ID: fromID},
	}
}

func TestJoinHandlerExecuteSendsCooldownMessage(t *testing.T) {
	w := &joinTestWorld{
		chat:  domain.Chat{ID: 100},
		users: map[int64]domain.User{1: {ID: 1, FirstName: "Wyatt"}},
		game: &domain.Game{
			ID:       uuid.Must(uuid.NewV7()),
			ChatID:   100,
			OwnerID:  1,
			PlayedAt: sql.NullTime{Time: time.Now(), Valid: true},
		},
	}
	w.chat.Settings.RequiredPlayers = 6

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, 1))

	msgs := bot.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "wait for game timeout", translator.Config{})[msgs[0]] {
		t.Fatalf("expected a cooldown message, got %q", msgs[0])
	}
}

func TestJoinHandlerExecuteSendsAlreadyInGameMessage(t *testing.T) {
	gameID := uuid.Must(uuid.NewV7())
	owner := domain.User{ID: 1, FirstName: "Wyatt"}
	w := &joinTestWorld{
		chat:  domain.Chat{ID: 100},
		users: map[int64]domain.User{1: owner},
		game: &domain.Game{
			ID:      gameID,
			ChatID:  100,
			OwnerID: 1,
			Owner:   owner,
		},
		gunslingers: []*domain.Gunslinger{
			{ID: uuid.Must(uuid.NewV7()), GameID: gameID, PlayerID: 1, Player: owner},
		},
	}
	w.chat.Settings.RequiredPlayers = 6

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, 1))

	msgs := bot.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "player already in game", translator.Config{})[msgs[0]] {
		t.Fatalf("expected an already-in-game message, got %q", msgs[0])
	}
}

func TestJoinHandlerExecuteSendsNotEnoughPlayersMessageToNonOwnerJoiner(t *testing.T) {
	gameID := uuid.Must(uuid.NewV7())
	owner := domain.User{ID: 1, FirstName: "Wyatt"}
	joiner := domain.User{ID: 2, FirstName: "Kate"}
	w := &joinTestWorld{
		chat:  domain.Chat{ID: 100},
		users: map[int64]domain.User{1: owner, 2: joiner},
		game: &domain.Game{
			ID:           gameID,
			ChatID:       100,
			OwnerID:      1,
			Owner:        owner,
			PlayersCount: 6,
		},
		gunslingers: []*domain.Gunslinger{
			{ID: uuid.Must(uuid.NewV7()), GameID: gameID, PlayerID: 1, Player: owner},
		},
	}
	w.chat.Settings.RequiredPlayers = 6

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, 2))

	msgs := bot.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "joining the game", translator.Config{})[msgs[0]] {
		t.Fatalf("expected a not-enough-players message, got %q", msgs[0])
	}
}

func TestJoinHandlerExecuteOwnerCreatingGameGetsOnlyCreationMessage(t *testing.T) {
	w := &joinTestWorld{
		chat:  domain.Chat{ID: 100},
		users: map[int64]domain.User{1: {ID: 1, FirstName: "Wyatt"}},
	}
	w.chat.Settings.RequiredPlayers = 6

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, 1))

	msgs := bot.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "game creation", translator.Config{})[msgs[0]] {
		t.Fatalf("expected a game-creation message, got %q", msgs[0])
	}
}

func twoPlayerGameAwaitingSecondJoiner() (*joinTestWorld, domain.User, domain.User) {
	gameID := uuid.Must(uuid.NewV7())
	owner := domain.User{ID: 1, FirstName: "Wyatt"}
	joiner := domain.User{ID: 2, FirstName: "Kate"}

	w := &joinTestWorld{
		chat:  domain.Chat{ID: 100},
		users: map[int64]domain.User{1: owner, 2: joiner},
		game: &domain.Game{
			ID:           gameID,
			ChatID:       100,
			OwnerID:      1,
			Owner:        owner,
			PlayersCount: 2,
		},
		gunslingers: []*domain.Gunslinger{
			{ID: uuid.Must(uuid.NewV7()), GameID: gameID, PlayerID: 1, Player: owner},
		},
	}
	w.chat.Settings.RequiredPlayers = 2

	return w, owner, joiner
}

func TestJoinHandlerExecuteAtomicBulletKillsEveryoneAndKicksAll(t *testing.T) {
	w, owner, joiner := twoPlayerGameAwaitingSecondJoiner()

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewAtomicBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, joiner.ID))

	msgs := bot.sentMessages()
	if len(msgs) != 3 {
		t.Fatalf("expected 2 'play the game' messages plus 1 outcome message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "killed by atomic bullet", translator.Config{})[msgs[2]] {
		t.Fatalf("expected the third message to announce the atomic bullet, got %q", msgs[2])
	}

	kicked := bot.kickedPlayers()
	if len(kicked) != 2 {
		t.Fatalf("expected both players to be kicked, got %v", kicked)
	}
	if kicked[0] != owner.ID || kicked[1] != joiner.ID {
		t.Fatalf("expected both the owner and the joiner to be kicked, got %v", kicked)
	}
}

func TestJoinHandlerExecuteLeadBulletKillsOneVictimAndNamesThem(t *testing.T) {
	w, _, joiner := twoPlayerGameAwaitingSecondJoiner()

	bot := &joinTestBot{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, joiner.ID))

	msgs := bot.sentMessages()
	if len(msgs) != 3 {
		t.Fatalf("expected 2 'play the game' messages plus 1 outcome message, got %d: %v", len(msgs), msgs)
	}

	kicked := bot.kickedPlayers()
	if len(kicked) != 1 {
		t.Fatalf("expected exactly one player to be kicked, got %v", kicked)
	}

	victim := w.users[kicked[0]]
	wantTexts := possibleTexts(trans, "gunslinger killed", translator.Config{
		Args: map[string]string{"%gunslinger": victim.Mention()},
	})
	if !wantTexts[msgs[2]] {
		t.Fatalf("expected the third message to name the victim %s, got %q", victim.Mention(), msgs[2])
	}
}

func TestJoinHandlerExecuteSendsMessageWhenKickFails(t *testing.T) {
	w, _, joiner := twoPlayerGameAwaitingSecondJoiner()

	bot := &joinTestBot{kickErr: errors.New("telegram: kick failed")}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	hdlr := newJoinTestHandler(w, bot, newTestNaganWithBullet(t, service.NewLeadBullet()))

	hdlr.Execute(context.Background(), joinMessage(100, joiner.ID))

	msgs := bot.sentMessages()
	if len(msgs) != 4 {
		t.Fatalf("expected 2 'play the game' messages, 1 kill message and 1 kick-failure message, got %d: %v", len(msgs), msgs)
	}
	if !possibleTexts(trans, "player is not kicked", translator.Config{})[msgs[3]] {
		t.Fatalf("expected the last message to report the failed kick, got %q", msgs[3])
	}
}
