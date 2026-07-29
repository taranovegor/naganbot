package service

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fakeHTTPClient struct {
	fail bool
}

func (c *fakeHTTPClient) Do(*http.Request) (*http.Response, error) {
	body := `{"ok":true,"result":{}}`
	if c.fail {
		body = `{"ok":false,"error_code":500,"description":"boom"}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
	}, nil
}

func newTestBot(t *testing.T, client *fakeHTTPClient) *Bot {
	t.Helper()

	botAPI, err := tgbotapi.NewBotAPIWithClient("test-token", tgbotapi.APIEndpoint, client)
	if err != nil {
		t.Fatalf("failed to build test bot: %v", err)
	}

	return NewBot(botAPI)
}

func TestSendMessageLogsOnAPIFailure(t *testing.T) {
	client := &fakeHTTPClient{}
	bot := newTestBot(t, client)

	var buf bytes.Buffer
	original := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(original) })

	client.fail = true
	bot.SendMessage(1, "hello")

	if !strings.Contains(buf.String(), "failed to send message to chat 1") {
		t.Fatalf("expected a log entry about the failed send, got: %q", buf.String())
	}
}

func TestBanReturnsErrorOnAPIFailure(t *testing.T) {
	client := &fakeHTTPClient{}
	bot := newTestBot(t, client)

	client.fail = true
	if err := bot.Ban(1, 2, 3); err == nil {
		t.Fatal("expected an error from Ban when the API call fails")
	}
}

func TestKickReturnsErrorOnAPIFailure(t *testing.T) {
	client := &fakeHTTPClient{}
	bot := newTestBot(t, client)

	client.fail = true
	if err := bot.Kick(1, 2); err == nil {
		t.Fatal("expected an error from Kick when the API call fails")
	}
}

func TestInlineKeyboardMarkupPreservesRowAndButtonOrder(t *testing.T) {
	keyboard := Keyboard{
		{{Data: "a1", Text: "one"}, {Data: "a2", Text: "two"}, {Data: "a3", Text: "three"}},
		{{Data: "b1", Text: "four"}},
	}

	got := inlineKeyboardMarkup(keyboard)

	if len(got.InlineKeyboard) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(got.InlineKeyboard))
	}

	firstRow := got.InlineKeyboard[0]
	if len(firstRow) != 3 {
		t.Fatalf("expected 3 buttons in the first row, got %d", len(firstRow))
	}

	wantFirstRow := []Button{{Data: "a1", Text: "one"}, {Data: "a2", Text: "two"}, {Data: "a3", Text: "three"}}
	for i, want := range wantFirstRow {
		got := firstRow[i]
		if got.Text != want.Text {
			t.Fatalf("button %d: expected text %q, got %q", i, want.Text, got.Text)
		}
		if got.CallbackData == nil || *got.CallbackData != want.Data {
			t.Fatalf("button %d: expected callback data %q, got %v", i, want.Data, got.CallbackData)
		}
	}

	secondRow := got.InlineKeyboard[1]
	if len(secondRow) != 1 {
		t.Fatalf("expected 1 button in the second row, got %d", len(secondRow))
	}
	if secondRow[0].Text != "four" || secondRow[0].CallbackData == nil || *secondRow[0].CallbackData != "b1" {
		t.Fatalf("unexpected second row button: %+v", secondRow[0])
	}
}

func TestInlineKeyboardMarkupDoesNotConfuseTextAndCallbackData(t *testing.T) {
	keyboard := Keyboard{{{Data: "required-players_6", Text: "6 patronov"}}}

	got := inlineKeyboardMarkup(keyboard)

	button := got.InlineKeyboard[0][0]
	if button.Text != "6 patronov" {
		t.Fatalf("expected button text %q, got %q", "6 patronov", button.Text)
	}
	if button.CallbackData == nil || *button.CallbackData != "required-players_6" {
		t.Fatalf("expected callback data %q, got %v", "required-players_6", button.CallbackData)
	}
}

func TestInlineKeyboardMarkupHandlesEmptyKeyboardWithoutPanicking(t *testing.T) {
	got := inlineKeyboardMarkup(nil)

	if len(got.InlineKeyboard) != 0 {
		t.Fatalf("expected no rows for an empty keyboard, got %d", len(got.InlineKeyboard))
	}
}
