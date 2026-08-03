package usecase

import (
	"testing"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

type lifecycleChatRepo struct {
	chat domain.Chat
}

func (r *lifecycleChatRepo) Get(int64) (domain.Chat, error)    { return r.chat, nil }
func (r *lifecycleChatRepo) Save(*domain.Chat) error           { return nil }
func (r *lifecycleChatRepo) UpdateSettings(*domain.Chat) error { return nil }

type lifecycleUserRepo struct {
	users map[int64]domain.User
}

func (r *lifecycleUserRepo) Get(id int64) (domain.User, error)       { return r.users[id], nil }
func (r *lifecycleUserRepo) GetByIDs([]int64) ([]domain.User, error) { return nil, nil }
func (r *lifecycleUserRepo) Save(*domain.User) error                 { return nil }

// lifecycleGunslingerStore is shared between the game and gunslinger fakes to mimic
// gorm's cascade-create of a Game's associated Gunslingers on Store — the exact
// mechanism the create+join ownership trap depends on: a gunslinger pre-populated
// into Game.Gunslingers gets persisted as a side effect of storing the game itself.
type lifecycleGunslingerStore struct {
	byGame map[uuid.UUID][]*domain.Gunslinger
}

func newLifecycleGunslingerStore() *lifecycleGunslingerStore {
	return &lifecycleGunslingerStore{byGame: make(map[uuid.UUID][]*domain.Gunslinger)}
}

func (s *lifecycleGunslingerStore) add(g *domain.Gunslinger) {
	s.byGame[g.GameID] = append(s.byGame[g.GameID], g)
}

func (s *lifecycleGunslingerStore) exists(userID int64, gameID uuid.UUID) bool {
	for _, g := range s.byGame[gameID] {
		if g.PlayerID == userID {
			return true
		}
	}
	return false
}

type lifecycleGameRepo struct {
	store *lifecycleGunslingerStore
}

func (r *lifecycleGameRepo) GetByID(uuid.UUID) (*domain.Game, error)                { return nil, nil }
func (r *lifecycleGameRepo) GetLatestGamesInChat(int64, int) ([]domain.Game, error) { return nil, nil }
func (r *lifecycleGameRepo) GetLatestForChat(int64) (*domain.Game, error)           { return nil, nil }
func (r *lifecycleGameRepo) GetActiveForChat(int64) (*domain.Game, error)           { return nil, nil }
func (r *lifecycleGameRepo) HasActiveOrCreatedTodayInChat(int64) bool               { return false }

func (r *lifecycleGameRepo) Store(game *domain.Game) error {
	for _, g := range game.Gunslingers {
		r.store.add(g)
	}
	return nil
}

type lifecycleGunslingerRepo struct {
	store *lifecycleGunslingerStore
}

func (r *lifecycleGunslingerRepo) Store(g *domain.Gunslinger) error {
	r.store.add(g)
	return nil
}

func (r *lifecycleGunslingerRepo) GetByGameID(uuid.UUID) ([]*domain.Gunslinger, error) {
	return nil, nil
}

func (r *lifecycleGunslingerRepo) IsPlayerExistsInGame(userID int64, gameID uuid.UUID) bool {
	return r.store.exists(userID, gameID)
}

func (r *lifecycleGunslingerRepo) GetTopShotPlayersInChat(int64) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}

func (r *lifecycleGunslingerRepo) GetTopShotPlayersByYearInChat(int64, int) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}

func (r *lifecycleGunslingerRepo) CountNumberOfPlayerGamesInChat(int64, int64) int64 { return 0 }
func (r *lifecycleGunslingerRepo) CountNumberOfSelfShotsInChat(int64, int64) int64   { return 0 }

func TestOwnerJoinsTheirOwnGameExactlyOnce(t *testing.T) {
	const (
		chatID  = int64(100)
		ownerID = int64(7)
	)

	owner := domain.User{ID: ownerID, FirstName: "Wyatt"}

	chat := domain.Chat{ID: chatID}
	chat.Settings.RequiredPlayers = 6

	store := newLifecycleGunslingerStore()
	chatRepo := &lifecycleChatRepo{chat: chat}
	userRepo := &lifecycleUserRepo{users: map[int64]domain.User{ownerID: owner}}
	gameRepo := &lifecycleGameRepo{store: store}
	gunslingerRepo := &lifecycleGunslingerRepo{store: store}

	createGame := NewCreateGameUseCase(service.NewLocker(), chatRepo, gameRepo, userRepo)
	joinGame := NewJoinGameUseCase(gameRepo, gunslingerRepo, userRepo)

	game, err := createGame.Execute(chatID, ownerID)
	if err != nil {
		t.Fatalf("unexpected error creating the game: %v", err)
	}
	if len(game.Gunslingers) != 0 {
		t.Fatalf("expected NewGame not to pre-populate Gunslingers, got %d", len(game.Gunslingers))
	}

	gunslinger, err := joinGame.Execute(game.ID, ownerID)
	if err != nil {
		t.Fatalf("expected the owner to be able to join the game they just created, got error: %v", err)
	}
	if gunslinger.PlayerID != ownerID {
		t.Fatalf("expected the gunslinger to belong to the owner, got player %d", gunslinger.PlayerID)
	}

	if !store.exists(ownerID, game.ID) {
		t.Fatal("expected the owner to be recorded as having joined the game")
	}
}
