package domain

import "testing"

func TestGame_MarkAsPlayed(t *testing.T) {
	var game *Game

	game = NewGame(0, 0, 6)
	game.MarkAsPlayed("lead", "https://example.com/proof")
	if game.IsPlayed() != true {
		t.Error()
	}
}
