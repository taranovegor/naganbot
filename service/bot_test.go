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
