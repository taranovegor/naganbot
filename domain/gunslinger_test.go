package domain

import "testing"

func TestGunslinger_MarkAsShotHimself_InNotPlayedGame(t *testing.T) {
	game := NewGame(0, 0, User{}, 6)
	gunslinger := NewGunslinger(game.ID, User{ID: 1})
	gunslinger.Game = game

	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}

	gunslinger.MarkAsShotHimself()
	if !gunslinger.ShotHimself {
		t.Error("gunslinger is expected to be shot")
	}
}

func TestGunslinger_MarkAsShotHimself_InPlayedGame(t *testing.T) {
	game := NewGame(0, 0, User{}, 6)
	game.MarkAsPlayed("lead", "https://example.com/proof")

	gunslinger := NewGunslinger(game.ID, User{ID: 1})
	gunslinger.Game = game

	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}

	gunslinger.MarkAsShotHimself()
	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}
}
