package domain

import (
	"testing"
	"time"
)

func TestNewDuel_LinksParticipantsAndIsNotPlayed(t *testing.T) {
	challenger := User{ID: 1}
	opponent := User{ID: 2}

	duel := NewDuel(10, challenger, opponent)

	if duel.ChatID != 10 || duel.ChallengerID != 1 || duel.OpponentID != 2 {
		t.Fatalf("unexpected duel: %+v", duel)
	}
	if duel.Challenger.ID != 1 || duel.Opponent.ID != 2 {
		t.Fatalf("participants are not attached: %+v", duel)
	}
	if duel.IsPlayed() {
		t.Fatal("a new duel must not be played")
	}
	if duel.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestNewDuel_GeneratesUniqueIDs(t *testing.T) {
	a := NewDuel(1, User{ID: 1}, User{ID: 2})
	b := NewDuel(1, User{ID: 1}, User{ID: 2})

	if a.ID == b.ID {
		t.Fatal("expected distinct duel ids")
	}
}

func TestDuel_MarkAsPlayed(t *testing.T) {
	duel := NewDuel(1, User{ID: 1}, User{ID: 2})

	duel.MarkAsPlayed(2, "https://example.com/proof")

	if !duel.IsPlayed() {
		t.Fatal("expected duel to be played")
	}
	if duel.LoserID != 2 || duel.ProofURL != "https://example.com/proof" {
		t.Fatalf("unexpected result: %+v", duel)
	}
}

func TestDuel_MarkAsPlayedIsIdempotent(t *testing.T) {
	duel := NewDuel(1, User{ID: 1}, User{ID: 2})
	duel.MarkAsPlayed(2, "first")
	playedAt := duel.PlayedAt

	duel.MarkAsPlayed(1, "second")

	if duel.LoserID != 2 || duel.ProofURL != "first" || duel.PlayedAt != playedAt {
		t.Fatalf("second MarkAsPlayed must not change the result: %+v", duel)
	}
}

func TestDuel_MarkAsDeclined(t *testing.T) {
	duel := NewDuel(1, User{ID: 1}, User{ID: 2})

	duel.MarkAsDeclined()

	if !duel.IsDeclined() || !duel.IsClosed() || duel.IsPlayed() {
		t.Fatalf("expected a declined, closed, unplayed duel: %+v", duel)
	}
}

func TestDuel_ClosedDuelCannotChangeState(t *testing.T) {
	declined := NewDuel(1, User{ID: 1}, User{ID: 2})
	declined.MarkAsDeclined()
	declined.MarkAsPlayed(1, "x")
	if declined.IsPlayed() {
		t.Fatal("a declined duel must not become played")
	}

	played := NewDuel(1, User{ID: 1}, User{ID: 2})
	played.MarkAsPlayed(1, "x")
	played.MarkAsDeclined()
	if played.IsDeclined() {
		t.Fatal("a played duel must not become declined")
	}
}

func TestDuel_IsExpired(t *testing.T) {
	duel := NewDuel(1, User{ID: 1}, User{ID: 2})

	if duel.IsExpired(duel.CreatedAt.Add(DuelTTL)) {
		t.Fatal("a duel must stay valid for the full TTL")
	}
	if !duel.IsExpired(duel.CreatedAt.Add(DuelTTL + time.Nanosecond)) {
		t.Fatal("a duel must expire after the TTL")
	}
}
