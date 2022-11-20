package game

func NewKnight(pos *Position, t *Team) *Knight {
	k := &Knight{}
	k.SetName("knight")
	k.Position = pos
	k.SetTeam(t)
	return k
}

// Knight is data type of chess figure
type Knight struct {
	Figure
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (k *Knight) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range k.detectionOfBrokenFields() {
		if !k.team.Figures.ExistsByPosition(position) && !k.kingOnTheBeatenFieldAfterMove(position) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (k *Knight) detectionOfBrokenFields() []*Position {
	var data []*Position

	pos := NewPosition(k.X+1, k.Y+2)
	if k.positionOnBoard(pos) {
		data = append(data, NewPosition(k.X+1, k.Y+2))
	}

	pos = NewPosition(k.X+2, k.Y+1)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X+2, k.Y-1)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X+1, k.Y-2)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X-1, k.Y-2)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X-2, k.Y-1)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X-2, k.Y+1)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	pos = NewPosition(k.X-1, k.Y+2)
	if k.positionOnBoard(pos) {
		data = append(data, pos)
	}

	return data
}

// Validation return true if this move are valid or return false
func (k *Knight) Validation(pos *Position) (bool, string) {
	if !k.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *k.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if k.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if k.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range k.detectionOfBrokenFields() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (k *Knight) Move(pos *Position) {
	k.team.pawnDoubleMove.clearPawnDoubleMove()
	k.MoveFigure(pos)
}
