package game

// Figurer a set of methods for any chess figure (king, queen, rook, knight, bishop, pawn)
type Figurer interface {
	GetName() string
	SetName(string)
	SetPosition(*Position)
	GetPosition() *Position
	Move(*Position)
	Validation(*Position) (bool, error)
	SetTeam(*Team)
	positionOnBoard(*Position) bool
	kingOnTheBeatenFieldAfterMove(*Position) bool
	GetPositionByDirectionAndRemote(Direction, uint8) *Position
	GetBrokenFields() *Positions
	GetPossibleMoves(thereIs bool) *Positions
	CanWalkLikeThat(*Position) bool
	IsAlreadyMove() bool
	setAlreadyMove(bool)
}
