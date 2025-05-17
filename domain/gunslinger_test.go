package domain

import "testing"

func TestGunslinger_MarkAsShotHimself_InNotPlayedGame(t *testing.T) {
	game := NewGame(0, 0)
	gunslinger := game.Gunslingers[0]

	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}

	gunslinger.MarkAsShotHimself()
	if !gunslinger.ShotHimself {
		t.Error("gunslinger is expected to be shot")
	}
}

func TestGunslinger_MarkAsShotHimself_InPlayedGame(t *testing.T) {
	game := NewGame(0, 0)
	game.MarkAsPlayed("lead")

	gunslinger := game.Gunslingers[0]
	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}

	gunslinger.MarkAsShotHimself()
	if gunslinger.ShotHimself {
		t.Error("gunslinger is not expected to be shot")
	}
}
