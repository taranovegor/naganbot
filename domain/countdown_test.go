package domain

import (
	"testing"
	"time"
)

func countdownGunslingers(joinedAt ...time.Time) []*Gunslinger {
	gunslingers := make([]*Gunslinger, 0, len(joinedAt))
	for i, at := range joinedAt {
		gunslinger := NewGunslinger(NewGame(1, 1, User{}, 6, GameModeCountdown).ID, User{ID: int64(i + 1)})
		gunslinger.JoinedAt = at
		gunslingers = append(gunslingers, gunslinger)
	}

	return gunslingers
}

func joinedAtSameTime(count int, at time.Time) []*Gunslinger {
	times := make([]time.Time, count)
	for i := range times {
		times[i] = at
	}

	return countdownGunslingers(times...)
}

func TestCountdownWindow(t *testing.T) {
	tests := map[int]time.Duration{
		0:  0,
		2:  0,
		3:  60 * time.Minute,
		4:  30 * time.Minute,
		5:  15 * time.Minute,
		10: 15 * time.Minute,
	}

	for players, expected := range tests {
		if got := CountdownWindow(players); got != expected {
			t.Errorf("CountdownWindow(%d) = %s, expected %s", players, got, expected)
		}
	}
}

func TestNextMidnight(t *testing.T) {
	noon := time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC)
	expected := time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)

	if got := NextMidnight(noon); !got.Equal(expected) {
		t.Errorf("expected %s, got %s", expected, got)
	}

	exactMidnight := time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC)
	if got := NextMidnight(exactMidnight); !got.Equal(expected) {
		t.Errorf("a time at midnight sharp must move to the next midnight, expected %s, got %s", expected, got)
	}
}

func TestGame_StartsAt(t *testing.T) {
	noon := time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		mode        GameMode
		gunslingers []*Gunslinger
		expected    time.Time
		expectedOK  bool
	}{
		{"fixed game has no start time", GameModeFixed, joinedAtSameTime(5, noon), time.Time{}, false},
		{"too few players", GameModeCountdown, joinedAtSameTime(2, noon), time.Time{}, false},
		{"three players wait an hour", GameModeCountdown, joinedAtSameTime(3, noon), noon.Add(time.Hour), true},
		{"four players wait half an hour", GameModeCountdown, joinedAtSameTime(4, noon), noon.Add(30 * time.Minute), true},
		{"five players wait a quarter", GameModeCountdown, joinedAtSameTime(5, noon), noon.Add(15 * time.Minute), true},
		{"full lobby starts at the last join", GameModeCountdown, joinedAtSameTime(10, noon), noon, true},
		{
			"the last join sets the clock",
			GameModeCountdown,
			countdownGunslingers(noon, noon, noon.Add(10*time.Minute)),
			noon.Add(70 * time.Minute),
			true,
		},
		{
			"midnight caps the wait",
			GameModeCountdown,
			joinedAtSameTime(3, time.Date(2026, time.March, 31, 23, 30, 0, 0, time.UTC)),
			time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(1, 1, User{}, 6, tt.mode)

			got, ok := game.StartsAt(tt.gunslingers)
			if ok != tt.expectedOK {
				t.Fatalf("expected ok=%t, got %t", tt.expectedOK, ok)
			}
			if !got.Equal(tt.expected) {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestGame_ShouldStart_Fixed(t *testing.T) {
	now := time.Now()
	game := NewGame(1, 1, User{}, 4, GameModeFixed)

	if game.ShouldStart(joinedAtSameTime(3, now), now) {
		t.Error("a fixed game must wait for every required player")
	}
	if !game.ShouldStart(joinedAtSameTime(4, now), now) {
		t.Error("a fixed game must start once it has the required players")
	}
	if !game.ShouldStart(joinedAtSameTime(4, now), now.Add(-time.Hour)) {
		t.Error("time must not matter to a fixed game")
	}
}

func TestGame_ShouldStart_Countdown(t *testing.T) {
	noon := time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC)
	game := NewGame(1, 1, User{}, 6, GameModeCountdown)

	if game.ShouldStart(joinedAtSameTime(2, noon), noon.Add(24*time.Hour)) {
		t.Error("a countdown game must never start without the minimum players")
	}
	if game.ShouldStart(joinedAtSameTime(3, noon), noon.Add(59*time.Minute)) {
		t.Error("a countdown game must wait for its deadline")
	}
	if !game.ShouldStart(joinedAtSameTime(3, noon), noon.Add(time.Hour)) {
		t.Error("a countdown game must start at its deadline")
	}
	if !game.ShouldStart(joinedAtSameTime(10, noon), noon) {
		t.Error("a full countdown game must start right away")
	}
}
