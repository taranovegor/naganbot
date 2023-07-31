package domain

import "testing"

func TestGame_MarkAsPlayed(t *testing.T) {
	var game *Game

	game = NewGame(0, 0)
	game.MarkAsPlayed()
	if game.IsPlayed() != true {
		t.Error()
	}
}
