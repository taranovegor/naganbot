package drand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestProofURLUsesTheConfiguredBaseURL(t *testing.T) {
	client := NewClientWithURL("https://relay.example.com")
	gameID := uuid.New()

	got := client.ProofURL(42, gameID)
	want := fmt.Sprintf("https://relay.example.com/%s/public/42?game[id]=%s", chainHash, gameID.String())

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestProofURLFallsBackToTheFirstRelayWithoutAConfiguredBaseURL(t *testing.T) {
	client := NewClient()
	gameID := uuid.New()

	got := client.ProofURL(42, gameID)
	want := fmt.Sprintf("%s/%s/public/42?game[id]=%s", relayURLs[0], chainHash, gameID.String())

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestGetLatestSucceedsOnAValidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Beacon{Round: 42, Randomness: "abc"})
	}))
	defer server.Close()

	beacon, err := NewClientWithURL(server.URL).GetLatest(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beacon.Round != 42 || beacon.Randomness != "abc" {
		t.Fatalf("unexpected beacon: %+v", beacon)
	}
}

func TestGetLatestReturnsAnErrorOnANon200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := NewClientWithURL(server.URL).GetLatest(context.Background())
	if err == nil {
		t.Fatal("expected an error for a non-200 response")
	}
}

func TestGetLatestReturnsAnErrorOnEmptyRandomness(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Beacon{Round: 1, Randomness: ""})
	}))
	defer server.Close()

	_, err := NewClientWithURL(server.URL).GetLatest(context.Background())
	if err == nil {
		t.Fatal("expected an error for empty randomness")
	}
}

func TestGetLatestFallsBackToTheNextRelayOnFailure(t *testing.T) {
	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failing.Close()

	working := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Beacon{Round: 7, Randomness: "ok"})
	}))
	defer working.Close()

	original := relayURLs
	relayURLs = []string{failing.URL, working.URL}
	t.Cleanup(func() { relayURLs = original })

	beacon, err := NewClient().GetLatest(context.Background())
	if err != nil {
		t.Fatalf("expected the second relay to succeed after the first failed, got error: %v", err)
	}
	if beacon.Randomness != "ok" {
		t.Fatalf("unexpected beacon: %+v", beacon)
	}
}

func TestGetLatestWrapsTheLastErrorWhenEveryRelayFails(t *testing.T) {
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer first.Close()

	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer second.Close()

	original := relayURLs
	relayURLs = []string{first.URL, second.URL}
	t.Cleanup(func() { relayURLs = original })

	_, err := NewClient().GetLatest(context.Background())
	if err == nil {
		t.Fatal("expected an error when every relay fails")
	}
	if !strings.Contains(err.Error(), "all drand relays failed") {
		t.Fatalf("expected the wrapped all-relays-failed message, got %v", err)
	}
	if !strings.Contains(err.Error(), "503") {
		t.Fatalf("expected the wrapped error to carry the last relay's failure (503), got %v", err)
	}
}
