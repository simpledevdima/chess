package game

// Rook is data type of chess figure
type Rook struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (rook *Rook) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range rook.detectionOfBrokenFields() {
		if !rook.team.FigureExist(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (rook *Rook) detectionOfBrokenFields() []Position {
	var data []Position
	directions := struct {
		top    bool
		right  bool
		bottom bool
		left   bool
	}{true, true, true, true}
	for i := 1; i <= 7; i++ {
		if directions.top && rook.coordsOnBoard(rook.X, rook.Y+i) {
			data = append(data, Position{X: rook.X, Y: rook.Y + i})
		}
		if directions.right && rook.coordsOnBoard(rook.X+i, rook.Y) {
			data = append(data, Position{X: rook.X + i, Y: rook.Y})
		}
		if directions.bottom && rook.coordsOnBoard(rook.X, rook.Y-i) {
			data = append(data, Position{X: rook.X, Y: rook.Y - i})
		}
		if directions.left && rook.coordsOnBoard(rook.X-i, rook.Y) {
			data = append(data, Position{X: rook.X - i, Y: rook.Y})
		}
		if rook.team.FigureExist(rook.X, rook.Y+i) ||
			rook.enemy.FigureExist(rook.X, rook.Y+i) ||
			!rook.coordsOnBoard(rook.X, rook.Y+i) {
			directions.top = false
		}
		if rook.team.FigureExist(rook.X+i, rook.Y) ||
			rook.enemy.FigureExist(rook.X+i, rook.Y) ||
			!rook.coordsOnBoard(rook.X+i, rook.Y) {
			directions.right = false
		}
		if rook.team.FigureExist(rook.X, rook.Y-i) ||
			rook.enemy.FigureExist(rook.X, rook.Y-i) ||
			!rook.coordsOnBoard(rook.X, rook.Y-i) {
			directions.bottom = false
		}
		if rook.team.FigureExist(rook.X-i, rook.Y) ||
			rook.enemy.FigureExist(rook.X-i, rook.Y) ||
			!rook.coordsOnBoard(rook.X-i, rook.Y) {
			directions.left = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (rook *Rook) Validation(x int, y int) (bool, string) {
	if !rook.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if rook.X == x && rook.Y == y {
		return false, "can't walk around"
	}
	if rook.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if rook.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if change x && y is not valid for rook
	if rook.X != x && rook.Y != y {
		return false, "rook doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range rook.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (rook *Rook) Move(x int, y int) {
	rook.team.pawnDoubleMove.clearPawnDoubleMove()
	rook.MoveFigure(x, y)
}
