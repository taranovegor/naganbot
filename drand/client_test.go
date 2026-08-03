package drand

import (
	"fmt"
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
