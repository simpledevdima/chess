package game

// Figurer a set of methods for any chess figure (king, queen, rook, knight, bishop, pawn)
type Figurer interface {
	GetName() string
	SetName(string)
	SetPosition(*Position)
	GetPosition() *Position
	Move(*Position)
	Validation(*Position) (bool, string)
	SetTeam(*Team)
	positionOnBoard(*Position) bool
	kingOnTheBeatenFieldAfterMove(*Position) bool
	GetBrokenFields() *Positions
	GetPossibleMoves() *Positions
	SetBrokenFields(*Positions)
	SetPossibleMoves(*Positions)
	IsAlreadyMove() bool
	setAlreadyMove(bool)
}
