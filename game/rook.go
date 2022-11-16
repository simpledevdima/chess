package game

func NewRook(x, y int, t *Team) *Rook {
	f := &Rook{}
	f.SetName("rook")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// Rook is data type of chess figure
type Rook struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (r *Rook) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range r.detectionOfBrokenFields() {
		if !r.team.FigureExist(position.X, position.Y) && !r.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (r *Rook) detectionOfBrokenFields() []*Position {
	var data []*Position
	directions := struct {
		top    bool
		right  bool
		bottom bool
		left   bool
	}{true, true, true, true}
	for i := 1; i <= 7; i++ {
		if directions.top && r.coordsOnBoard(r.X, r.Y+i) {
			data = append(data, NewPosition(r.X, r.Y+i))
		}
		if directions.right && r.coordsOnBoard(r.X+i, r.Y) {
			data = append(data, NewPosition(r.X+i, r.Y))
		}
		if directions.bottom && r.coordsOnBoard(r.X, r.Y-i) {
			data = append(data, NewPosition(r.X, r.Y-i))
		}
		if directions.left && r.coordsOnBoard(r.X-i, r.Y) {
			data = append(data, NewPosition(r.X-i, r.Y))
		}
		if r.team.FigureExist(r.X, r.Y+i) ||
			r.team.enemy.FigureExist(r.X, r.Y+i) ||
			!r.coordsOnBoard(r.X, r.Y+i) {
			directions.top = false
		}
		if r.team.FigureExist(r.X+i, r.Y) ||
			r.team.enemy.FigureExist(r.X+i, r.Y) ||
			!r.coordsOnBoard(r.X+i, r.Y) {
			directions.right = false
		}
		if r.team.FigureExist(r.X, r.Y-i) ||
			r.team.enemy.FigureExist(r.X, r.Y-i) ||
			!r.coordsOnBoard(r.X, r.Y-i) {
			directions.bottom = false
		}
		if r.team.FigureExist(r.X-i, r.Y) ||
			r.team.enemy.FigureExist(r.X-i, r.Y) ||
			!r.coordsOnBoard(r.X-i, r.Y) {
			directions.left = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (r *Rook) Validation(x int, y int) (bool, string) {
	if !r.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if r.X == x && r.Y == y {
		return false, "can't walk around"
	}
	if r.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if r.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if change x && y is not valid for rook
	if r.X != x && r.Y != y {
		return false, "rook doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range r.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (r *Rook) Move(x int, y int) {
	r.team.pawnDoubleMove.clearPawnDoubleMove()
	r.MoveFigure(x, y)
}
