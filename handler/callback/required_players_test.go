package callback

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

func TestRevolverKeyboardListsAllOptionsInOrderWithSelectedMarked(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	keyboard := RevolverKeyboard(6, trans)

	if len(keyboard) != len(revolverOptions) {
		t.Fatalf("expected %d rows, got %d", len(revolverOptions), len(keyboard))
	}

	for i, n := range revolverOptions {
		row := keyboard[i]
		if len(row) != 1 {
			t.Fatalf("row %d: expected exactly one button, got %d", i, len(row))
		}

		button := row[0]

		wantData := RequiredPlayers.SetArgs(fmt.Sprint(n)).ToString()
		if button.Data != wantData {
			t.Fatalf("row %d: expected callback data %q, got %q", i, wantData, button.Data)
		}

		wantText := trans.Get(fmt.Sprintf("%d shot revolver", n), translator.Config{})
		if n == 6 {
			wantText = fmt.Sprintf("🔫 %s", wantText)
		}
		if button.Text != wantText {
			t.Fatalf("row %d: expected text %q, got %q", i, wantText, button.Text)
		}
	}
}

func TestRevolverKeyboardMarksOnlyTheSelectedOption(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	keyboard := RevolverKeyboard(4, trans)

	marked := 0
	for i, n := range revolverOptions {
		text := keyboard[i][0].Text
		if strings.HasPrefix(text, "🔫") {
			marked++
			if n != 4 {
				t.Fatalf("expected only option 4 to be marked, but option %d was marked too", n)
			}
		}
	}
	if marked != 1 {
		t.Fatalf("expected exactly one marked option, got %d", marked)
	}
}

type fakeMessenger struct {
	isAdmin    bool
	isAdminErr error

	answers        []string
	editedKeyboard *service.Keyboard
}

func (b *fakeMessenger) AnswerCallback(_ string, text string) {
	b.answers = append(b.answers, text)
}
func (b *fakeMessenger) EditMessageReplyMarkup(_ int64, _ int, keyboard service.Keyboard) {
	b.editedKeyboard = &keyboard
}
func (b *fakeMessenger) IsAdmin(int64, int64) (bool, error) {
	return b.isAdmin, b.isAdminErr
}

type fakeCallbackChatRepo struct {
	chat      domain.Chat
	getErr    error
	updateErr error
	updated   *domain.Chat
}

func (r *fakeCallbackChatRepo) Get(int64) (domain.Chat, error) { return r.chat, r.getErr }
func (r *fakeCallbackChatRepo) Save(*domain.Chat) error        { return nil }
func (r *fakeCallbackChatRepo) UpdateSettings(chat *domain.Chat) error {
	r.updated = chat
	return r.updateErr
}

func requiredPlayersQuery(chatID int64, fromID int64, messageID int, players string) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		ID:      "cb1",
		Data:    RequiredPlayers.SetArgs(players).ToString(),
		From:    &tgbotapi.User{ID: fromID},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: chatID}, MessageID: messageID},
	}
}

func somethingWentWrongTexts(trans *translator.Translator) map[string]bool {
	texts := make(map[string]bool)
	for i := 0; i < 50; i++ {
		texts[trans.Get("something went wrong", translator.Config{})] = true
	}
	return texts
}

func TestRequiredPlayersExecuteSendsSomethingWentWrongWhenIsAdminCheckFails(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdminErr: errors.New("telegram: request failed")}
	hdlr := &requiredPlayers{chatRepo: &fakeCallbackChatRepo{}, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), requiredPlayersQuery(100, 7, 1, "6"))

	if len(bot.answers) != 1 || !somethingWentWrongTexts(trans)[bot.answers[0]] {
		t.Fatalf("expected a single 'something went wrong' answer, got %v", bot.answers)
	}
	if bot.editedKeyboard != nil {
		t.Fatal("expected no keyboard edit when the admin check fails")
	}
}

func TestRequiredPlayersExecuteNotifiesNonAdmins(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: false}
	chatRepo := &fakeCallbackChatRepo{}
	hdlr := &requiredPlayers{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), requiredPlayersQuery(100, 7, 1, "6"))

	wantTexts := make(map[string]bool)
	for i := 0; i < 50; i++ {
		wantTexts[trans.Get("settings can be changed only by admins", translator.Config{})] = true
	}
	if len(bot.answers) != 1 || !wantTexts[bot.answers[0]] {
		t.Fatalf("expected a single admins-only answer, got %v", bot.answers)
	}
	if chatRepo.updated != nil {
		t.Fatal("expected settings not to be touched for a non-admin")
	}
}

func TestRequiredPlayersExecuteSendsSomethingWentWrongOnABadArgument(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	hdlr := &requiredPlayers{chatRepo: &fakeCallbackChatRepo{}, bot: bot, trans: trans}

	query := requiredPlayersQuery(100, 7, 1, "6")
	query.Data = RequiredPlayers.ToString() + "_not-a-number"

	hdlr.Execute(context.Background(), query)

	if len(bot.answers) != 1 || !somethingWentWrongTexts(trans)[bot.answers[0]] {
		t.Fatalf("expected a single 'something went wrong' answer for a bad argument, got %v", bot.answers)
	}
}

func TestRequiredPlayersExecuteSendsSomethingWentWrongWhenChatLookupFails(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{getErr: errors.New("connection refused")}
	hdlr := &requiredPlayers{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), requiredPlayersQuery(100, 7, 1, "6"))

	if len(bot.answers) != 1 || !somethingWentWrongTexts(trans)[bot.answers[0]] {
		t.Fatalf("expected a single 'something went wrong' answer when the chat lookup fails, got %v", bot.answers)
	}
}

func TestRequiredPlayersExecuteSendsSomethingWentWrongWhenUpdateSettingsFails(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{updateErr: errors.New("write failed")}
	hdlr := &requiredPlayers{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), requiredPlayersQuery(100, 7, 1, "6"))

	if len(bot.answers) != 1 || !somethingWentWrongTexts(trans)[bot.answers[0]] {
		t.Fatalf("expected a single 'something went wrong' answer when UpdateSettings fails, got %v", bot.answers)
	}
	if bot.editedKeyboard != nil {
		t.Fatal("expected no keyboard edit when UpdateSettings fails")
	}
}

func TestRequiredPlayersExecuteUpdatesSettingsAndEditsTheKeyboardOnSuccess(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{chat: domain.Chat{ID: 100}}
	hdlr := &requiredPlayers{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), requiredPlayersQuery(100, 7, 55, "8"))

	if chatRepo.updated == nil || chatRepo.updated.Settings.RequiredPlayers != 8 {
		t.Fatalf("expected settings to be updated with 8 required players, got %+v", chatRepo.updated)
	}

	wantAnswer := fmt.Sprintf(
		"%s\n%s",
		trans.Get("revolver has been replaced", translator.Config{Count: 8}),
		trans.Get("settings will be applied for next games", translator.Config{}),
	)
	if len(bot.answers) != 1 || bot.answers[0] != wantAnswer {
		t.Fatalf("expected answer %q, got %v", wantAnswer, bot.answers)
	}

	if bot.editedKeyboard == nil {
		t.Fatal("expected the keyboard to be edited on success")
	}
	wantKeyboard := RevolverKeyboard(8, trans)
	if !reflect.DeepEqual(*bot.editedKeyboard, wantKeyboard) {
		t.Fatalf("expected keyboard %+v, got %+v", wantKeyboard, *bot.editedKeyboard)
	}
}
