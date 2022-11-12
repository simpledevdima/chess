package game

// Knight is data type of chess figure
type Knight struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (knight *Knight) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range knight.detectionOfBrokenFields() {
		if !knight.team.FigureExist(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (knight *Knight) detectionOfBrokenFields() []Position {
	var data []Position

	if knight.coordsOnBoard(knight.X+1, knight.Y+2) {
		data = append(data, Position{X: knight.X + 1, Y: knight.Y + 2})
	}
	if knight.coordsOnBoard(knight.X+2, knight.Y+1) {
		data = append(data, Position{X: knight.X + 2, Y: knight.Y + 1})
	}
	if knight.coordsOnBoard(knight.X+2, knight.Y-1) {
		data = append(data, Position{X: knight.X + 2, Y: knight.Y - 1})
	}
	if knight.coordsOnBoard(knight.X+1, knight.Y-2) {
		data = append(data, Position{X: knight.X + 1, Y: knight.Y - 2})
	}
	if knight.coordsOnBoard(knight.X-1, knight.Y-2) {
		data = append(data, Position{X: knight.X - 1, Y: knight.Y - 2})
	}
	if knight.coordsOnBoard(knight.X-2, knight.Y-1) {
		data = append(data, Position{X: knight.X - 2, Y: knight.Y - 1})
	}
	if knight.coordsOnBoard(knight.X-2, knight.Y+1) {
		data = append(data, Position{X: knight.X - 2, Y: knight.Y + 1})
	}
	if knight.coordsOnBoard(knight.X-1, knight.Y+2) {
		data = append(data, Position{X: knight.X - 1, Y: knight.Y + 2})
	}

	return data
}

// Validation return true if this move are valid or return false
func (knight *Knight) Validation(x int, y int) (bool, string) {
	if !knight.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if knight.X == x && knight.Y == y {
		return false, "can't walk around"
	}
	if knight.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if knight.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range knight.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (knight *Knight) Move(x int, y int) {
	knight.team.pawnDoubleMove.clearPawnDoubleMove()
	knight.MoveFigure(x, y)
}
