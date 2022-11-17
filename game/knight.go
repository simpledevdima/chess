package game

func NewKnight(x, y int, t *Team) *Knight {
	f := &Knight{}
	f.SetName("knight")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// Knight is data type of chess figure
type Knight struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (k *Knight) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range k.detectionOfBrokenFields() {
		if !k.team.Figures.ExistsByCoords(position.X, position.Y) && !k.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (k *Knight) detectionOfBrokenFields() []*Position {
	var data []*Position

	if k.coordsOnBoard(k.X+1, k.Y+2) {
		data = append(data, NewPosition(k.X+1, k.Y+2))
	}
	if k.coordsOnBoard(k.X+2, k.Y+1) {
		data = append(data, NewPosition(k.X+2, k.Y+1))
	}
	if k.coordsOnBoard(k.X+2, k.Y-1) {
		data = append(data, NewPosition(k.X+2, k.Y-1))
	}
	if k.coordsOnBoard(k.X+1, k.Y-2) {
		data = append(data, NewPosition(k.X+1, k.Y-2))
	}
	if k.coordsOnBoard(k.X-1, k.Y-2) {
		data = append(data, NewPosition(k.X-1, k.Y-2))
	}
	if k.coordsOnBoard(k.X-2, k.Y-1) {
		data = append(data, NewPosition(k.X-2, k.Y-1))
	}
	if k.coordsOnBoard(k.X-2, k.Y+1) {
		data = append(data, NewPosition(k.X-2, k.Y+1))
	}
	if k.coordsOnBoard(k.X-1, k.Y+2) {
		data = append(data, NewPosition(k.X-1, k.Y+2))
	}

	return data
}

// Validation return true if this move are valid or return false
func (k *Knight) Validation(x int, y int) (bool, string) {
	if !k.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if k.X == x && k.Y == y {
		return false, "can't walk around"
	}
	if k.team.Figures.ExistsByCoords(x, y) {
		return false, "this place is occupied by your figure"
	}
	if k.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range k.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (k *Knight) Move(x int, y int) {
	k.team.pawnDoubleMove.clearPawnDoubleMove()
	k.MoveFigure(x, y)
}
