package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type fakeHTTPClient struct {
	sentTexts []string
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		c.sentTexts = append(c.sentTexts, string(body))
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"ok":true,"result":{}}`))),
	}, nil
}

func newTestBot(t *testing.T) (*service.Bot, *fakeHTTPClient) {
	t.Helper()

	client := &fakeHTTPClient{}
	botAPI, err := tgbotapi.NewBotAPIWithClient("test-token", tgbotapi.APIEndpoint, client)
	if err != nil {
		t.Fatalf("failed to build test bot: %v", err)
	}

	return service.NewBot(botAPI), client
}

type fakeChatRepo struct {
	mu    sync.Mutex
	saved []*domain.Chat
}

func (r *fakeChatRepo) Get(int64) (domain.Chat, error)    { return domain.Chat{}, nil }
func (r *fakeChatRepo) UpdateSettings(*domain.Chat) error { return nil }
func (r *fakeChatRepo) Save(chat *domain.Chat) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.saved = append(r.saved, chat)
	return nil
}

type fakeUserRepo struct {
	mu    sync.Mutex
	saved []*domain.User
}

func (r *fakeUserRepo) Get(int64) (domain.User, error)          { return domain.User{}, nil }
func (r *fakeUserRepo) GetByIDs([]int64) ([]domain.User, error) { return nil, nil }
func (r *fakeUserRepo) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.saved = append(r.saved, user)
	return nil
}

type fakeCommandHandler struct {
	name string

	mu          sync.Mutex
	executed    bool
	capturedCtx context.Context
	onExecute   func()
}

func (h *fakeCommandHandler) Name() string { return h.name }
func (h *fakeCommandHandler) Execute(ctx context.Context, _ *tgbotapi.Message) {
	h.mu.Lock()
	h.executed = true
	h.capturedCtx = ctx
	onExecute := h.onExecute
	h.mu.Unlock()

	if onExecute != nil {
		onExecute()
	}
}

type fakeCallbackHandler struct {
	pattern callback.Pattern

	mu       sync.Mutex
	executed bool
}

func (h *fakeCallbackHandler) Pattern() callback.Pattern { return h.pattern }
func (h *fakeCallbackHandler) Execute(context.Context, *tgbotapi.CallbackQuery) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.executed = true
}

func groupChat() *tgbotapi.Chat {
	return &tgbotapi.Chat{ID: 100, Type: "group", Title: "Saloon"}
}

func testUser() *tgbotapi.User {
	return &tgbotapi.User{ID: 7, FirstName: "Wyatt", UserName: "wyatt"}
}

func commandMessage(command string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat:     groupChat(),
		From:     testUser(),
		Text:     command,
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command)}},
	}
}

func TestHandleUpdateSavesChatAndUserBeforeDispatchingCommand(t *testing.T) {
	chatRepo := &fakeChatRepo{}
	userRepo := &fakeUserRepo{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	savedBeforeExecute := false
	handler := &fakeCommandHandler{name: "join", onExecute: func() {
		savedBeforeExecute = len(chatRepo.saved) == 1 && len(userRepo.saved) == 1
	}}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      chatRepo,
		Users:      userRepo,
		Commands:   command.NewRegistry("", handler),
		Callbacks:  callback.NewRegistry(),
	}

	update := tgbotapi.Update{Message: commandMessage("/join")}
	a.HandleUpdate(context.Background(), update)

	if !handler.executed {
		t.Fatal("expected command handler to be executed")
	}
	if !savedBeforeExecute {
		t.Fatal("expected chat and user to be saved before command dispatch")
	}
	if len(chatRepo.saved) != 1 || chatRepo.saved[0].ID != 100 {
		t.Fatalf("unexpected saved chat: %+v", chatRepo.saved)
	}
	if len(userRepo.saved) != 1 || userRepo.saved[0].ID != 7 {
		t.Fatalf("unexpected saved user: %+v", userRepo.saved)
	}
}

func TestHandleUpdateSkipsSyncAndDispatchForPrivateChat(t *testing.T) {
	chatRepo := &fakeChatRepo{}
	userRepo := &fakeUserRepo{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, client := newTestBot(t)

	handler := &fakeCommandHandler{name: "join"}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      chatRepo,
		Users:      userRepo,
		Commands:   command.NewRegistry("", handler),
		Callbacks:  callback.NewRegistry(),
	}

	privateMessage := commandMessage("/join")
	privateMessage.Chat = &tgbotapi.Chat{ID: 7, Type: "private"}

	update := tgbotapi.Update{Message: privateMessage}
	a.HandleUpdate(context.Background(), update)

	if handler.executed {
		t.Fatal("command must not run for a private chat")
	}
	if len(chatRepo.saved) != 0 || len(userRepo.saved) != 0 {
		t.Fatal("chat and user must not be synced for a private chat")
	}
	if len(client.sentTexts) == 0 {
		t.Fatal("expected the bot to notify the user that the game is chat-only")
	}
}

func TestHandleUpdateSyncsButDoesNotDispatchPlainMessage(t *testing.T) {
	chatRepo := &fakeChatRepo{}
	userRepo := &fakeUserRepo{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      chatRepo,
		Users:      userRepo,
		Commands:   command.NewRegistry(""),
		Callbacks:  callback.NewRegistry(),
	}

	plainMessage := &tgbotapi.Message{Chat: groupChat(), From: testUser(), Text: "hello"}
	update := tgbotapi.Update{Message: plainMessage}
	a.HandleUpdate(context.Background(), update)

	if len(chatRepo.saved) != 1 || len(userRepo.saved) != 1 {
		t.Fatal("expected chat and user to still be synced for a plain message")
	}
}

func TestHandleUpdateDispatchesCallbackQueryAfterSync(t *testing.T) {
	chatRepo := &fakeChatRepo{}
	userRepo := &fakeUserRepo{}
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	handler := &fakeCallbackHandler{pattern: callback.RequiredPlayers}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      chatRepo,
		Users:      userRepo,
		Commands:   command.NewRegistry(""),
		Callbacks:  callback.NewRegistry(handler),
	}

	update := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Data:    callback.RequiredPlayers.ToString(),
			From:    testUser(),
			Message: &tgbotapi.Message{Chat: groupChat()},
		},
	}
	a.HandleUpdate(context.Background(), update)

	if !handler.executed {
		t.Fatal("expected callback handler to be executed")
	}
	if len(chatRepo.saved) != 1 || len(userRepo.saved) != 1 {
		t.Fatal("expected chat and user to be synced before callback dispatch")
	}
}

func TestHandleUpdatePropagatesTheCallerContextToTheCommandHandler(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot, _ := newTestBot(t)

	handler := &fakeCommandHandler{name: "join"}

	a := &App{
		Translator: trans,
		Bot:        bot,
		Chats:      &fakeChatRepo{},
		Users:      &fakeUserRepo{},
		Commands:   command.NewRegistry("", handler),
		Callbacks:  callback.NewRegistry(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	update := tgbotapi.Update{Message: commandMessage("/join")}
	a.HandleUpdate(ctx, update)

	if !handler.executed {
		t.Fatal("expected the command handler to be executed")
	}
	if !errors.Is(handler.capturedCtx.Err(), context.Canceled) {
		t.Fatalf("expected the handler to receive the caller's cancelled context, got err %v", handler.capturedCtx.Err())
	}
}
