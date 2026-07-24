package domain

type GameplayUnitOfWork interface {
	CommitPlayedGame(game *Game, victims []*Gunslinger) error
}
