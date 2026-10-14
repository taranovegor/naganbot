package domain

import "time"

const (
	CountdownMinPlayers = 3
	CountdownMaxPlayers = 10

	// LobbyIdleTTL is how long a lobby can stay without new players before it
	// stops blocking its chat.
	LobbyIdleTTL = 6 * time.Hour
)

// CountdownWindow is how long a countdown lobby waits for more players.
// The more players it has, the shorter the wait.
func CountdownWindow(players int) time.Duration {
	switch {
	case players < CountdownMinPlayers:
		return 0
	case players == CountdownMinPlayers:
		return 60 * time.Minute
	case players == CountdownMinPlayers+1:
		return 30 * time.Minute
	default:
		return 15 * time.Minute
	}
}

// NextMidnight returns the first midnight after t, in the same location as t.
func NextMidnight(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d+1, 0, 0, 0, 0, t.Location())
}

// StartsAt returns the time when the game starts by itself. It is calculated
// from the join times, so nothing is stored and a restart changes nothing.
// The second value is false when there is no start time: the game is fixed,
// or a countdown game has fewer than CountdownMinPlayers players.
func (g *Game) StartsAt(gunslingers []*Gunslinger) (time.Time, bool) {
	if g.Mode != GameModeCountdown || len(gunslingers) < CountdownMinPlayers {
		return time.Time{}, false
	}

	var lastJoin time.Time
	for _, gunslinger := range gunslingers {
		if gunslinger.JoinedAt.After(lastJoin) {
			lastJoin = gunslinger.JoinedAt
		}
	}

	if len(gunslingers) >= CountdownMaxPlayers {
		return lastJoin, true
	}

	startsAt := lastJoin.Add(CountdownWindow(len(gunslingers)))
	if midnight := NextMidnight(lastJoin); midnight.Before(startsAt) {
		startsAt = midnight
	}

	return startsAt, true
}

// ShouldStart reports whether the game can be played at the time now.
func (g *Game) ShouldStart(gunslingers []*Gunslinger, now time.Time) bool {
	if g.Mode != GameModeCountdown {
		return len(gunslingers) >= g.PlayersCount
	}

	startsAt, ok := g.StartsAt(gunslingers)

	return ok && !now.Before(startsAt)
}
