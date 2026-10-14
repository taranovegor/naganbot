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
	"github.com/taranovegor/naganbot/translator"
)

func gameModeQuery(chatID int64, fromID int64, messageID int, mode string) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		ID:      "cb1",
		Data:    GameMode.SetArgs(mode).ToString(),
		From:    &tgbotapi.User{ID: fromID},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: chatID}, MessageID: messageID},
	}
}

func TestSettingsKeyboardMarksTheSelectedModeAndOffersRevolversOnlyForFixed(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	fixed := SettingsKeyboard(domain.GameModeFixed, 6, trans)
	if len(fixed) != len(gameModes)+len(revolverOptions) {
		t.Fatalf("expected the modes followed by the revolvers, got %d rows", len(fixed))
	}
	if !strings.HasPrefix(fixed[0][0].Text, "✅") || strings.HasPrefix(fixed[1][0].Text, "✅") {
		t.Fatalf("expected only the fixed mode to be marked, got %q and %q", fixed[0][0].Text, fixed[1][0].Text)
	}
	if fixed[0][0].Data != GameMode.SetArgs("fixed").ToString() || fixed[1][0].Data != GameMode.SetArgs("countdown").ToString() {
		t.Fatalf("unexpected callback data: %q, %q", fixed[0][0].Data, fixed[1][0].Data)
	}

	countdown := SettingsKeyboard(domain.GameModeCountdown, 6, trans)
	if len(countdown) != len(gameModes) {
		t.Fatalf("expected only the modes for a countdown chat, got %d rows", len(countdown))
	}
	if strings.HasPrefix(countdown[0][0].Text, "✅") || !strings.HasPrefix(countdown[1][0].Text, "✅") {
		t.Fatalf("expected only the countdown mode to be marked, got %q and %q", countdown[0][0].Text, countdown[1][0].Text)
	}
}

func TestSettingsKeyboardTreatsAnUnsetModeAsFixed(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	if !reflect.DeepEqual(SettingsKeyboard("", 6, trans), SettingsKeyboard(domain.GameModeFixed, 6, trans)) {
		t.Fatal("expected an unset mode to look like the fixed one")
	}
}

func TestGameModeExecuteSendsSomethingWentWrongWhenIsAdminCheckFails(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdminErr: errors.New("telegram: request failed")}
	hdlr := &gameMode{chatRepo: &fakeCallbackChatRepo{}, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), gameModeQuery(100, 7, 1, "countdown"))

	if len(bot.answers) != 1 || bot.answers[0] != somethingWentWrongText(trans) {
		t.Fatalf("expected a single 'something went wrong' answer, got %v", bot.answers)
	}
}

func TestGameModeExecuteNotifiesNonAdminsAndChangesNothing(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: false}
	chatRepo := &fakeCallbackChatRepo{}
	hdlr := &gameMode{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), gameModeQuery(100, 7, 1, "countdown"))

	want := trans.Get("settings can be changed only by admins", translator.Config{})
	if len(bot.answers) != 1 || bot.answers[0] != want {
		t.Fatalf("expected the admins-only answer, got %v", bot.answers)
	}
	if chatRepo.updated != nil || bot.editedKeyboard != nil {
		t.Fatal("expected nothing to change for a non-admin")
	}
}

func TestGameModeExecuteRejectsAnUnknownMode(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{}
	hdlr := &gameMode{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), gameModeQuery(100, 7, 1, "anarchy"))

	if len(bot.answers) != 1 || bot.answers[0] != somethingWentWrongText(trans) {
		t.Fatalf("expected 'something went wrong', got %v", bot.answers)
	}
	if chatRepo.updated != nil {
		t.Fatal("expected an unknown mode never to be stored")
	}
}

func TestGameModeExecuteSendsSomethingWentWrongWhenTheChatCannotBeLoadedOrSaved(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	for name, repo := range map[string]*fakeCallbackChatRepo{
		"get":    {getErr: errors.New("db is down")},
		"update": {updateErr: errors.New("db is down")},
	} {
		t.Run(name, func(t *testing.T) {
			bot := &fakeMessenger{isAdmin: true}
			hdlr := &gameMode{chatRepo: repo, bot: bot, trans: trans}

			hdlr.Execute(context.Background(), gameModeQuery(100, 7, 1, "countdown"))

			if len(bot.answers) != 1 || bot.answers[0] != somethingWentWrongText(trans) {
				t.Fatalf("expected 'something went wrong', got %v", bot.answers)
			}
			if bot.editedKeyboard != nil {
				t.Fatal("expected no keyboard edit on failure")
			}
		})
	}
}

func TestGameModeExecuteDoesNothingWhenTheModeIsAlreadySet(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{chat: domain.Chat{ID: 100}}
	chatRepo.chat.Settings.Mode = domain.GameModeCountdown
	hdlr := &gameMode{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), gameModeQuery(100, 7, 1, "countdown"))

	if chatRepo.updated != nil || bot.editedKeyboard != nil {
		t.Fatal("expected no write and no keyboard edit for an unchanged mode")
	}
	if len(bot.answers) != 1 {
		t.Fatalf("expected the callback to be answered once, got %v", bot.answers)
	}
}

func TestGameModeExecuteStoresTheModeAndRedrawsTheKeyboard(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	bot := &fakeMessenger{isAdmin: true}
	chatRepo := &fakeCallbackChatRepo{chat: domain.Chat{ID: 100}}
	chatRepo.chat.Settings.RequiredPlayers = 7
	hdlr := &gameMode{chatRepo: chatRepo, bot: bot, trans: trans}

	hdlr.Execute(context.Background(), gameModeQuery(100, 7, 55, "countdown"))

	if chatRepo.updated == nil || chatRepo.updated.Settings.Mode != domain.GameModeCountdown {
		t.Fatalf("expected the countdown mode to be stored, got %+v", chatRepo.updated)
	}

	wantAnswer := fmt.Sprintf(
		"%s\n%s",
		trans.Get("countdown mode has been set", translator.Config{}),
		trans.Get("settings will be applied for next games", translator.Config{}),
	)
	if len(bot.answers) != 1 || bot.answers[0] != wantAnswer {
		t.Fatalf("expected answer %q, got %v", wantAnswer, bot.answers)
	}

	wantKeyboard := SettingsKeyboard(domain.GameModeCountdown, 7, trans)
	if bot.editedKeyboard == nil || !reflect.DeepEqual(*bot.editedKeyboard, wantKeyboard) {
		t.Fatalf("expected the keyboard to be redrawn as %+v, got %+v", wantKeyboard, bot.editedKeyboard)
	}
}
