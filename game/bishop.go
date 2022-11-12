package game

// Bishop is data type of chess figure
type Bishop struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (bishop *Bishop) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range bishop.detectionOfBrokenFields() {
		if !bishop.team.FigureExist(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (bishop *Bishop) detectionOfBrokenFields() []Position {
	var data []Position
	directions := struct {
		rightTop    bool
		rightBottom bool
		leftBottom  bool
		leftTop     bool
	}{true, true, true, true}
	for i := 1; i <= 7; i++ {
		if directions.rightTop && bishop.coordsOnBoard(bishop.X+i, bishop.Y+i) {
			data = append(data, Position{X: bishop.X + i, Y: bishop.Y + i})
		}
		if directions.rightBottom && bishop.coordsOnBoard(bishop.X+i, bishop.Y-i) {
			data = append(data, Position{X: bishop.X + i, Y: bishop.Y - i})
		}
		if directions.leftBottom && bishop.coordsOnBoard(bishop.X-i, bishop.Y-i) {
			data = append(data, Position{X: bishop.X - i, Y: bishop.Y - i})
		}
		if directions.leftTop && bishop.coordsOnBoard(bishop.X-i, bishop.Y+i) {
			data = append(data, Position{X: bishop.X - i, Y: bishop.Y + i})
		}
		if bishop.team.FigureExist(bishop.X+i, bishop.Y+i) ||
			bishop.enemy.FigureExist(bishop.X+i, bishop.Y+i) ||
			!bishop.coordsOnBoard(bishop.X+i, bishop.Y+i) {
			directions.rightTop = false
		}
		if bishop.team.FigureExist(bishop.X+i, bishop.Y-i) ||
			bishop.enemy.FigureExist(bishop.X+i, bishop.Y-i) ||
			!bishop.coordsOnBoard(bishop.X+i, bishop.Y-i) {
			directions.rightBottom = false
		}
		if bishop.team.FigureExist(bishop.X-i, bishop.Y-i) ||
			bishop.enemy.FigureExist(bishop.X-i, bishop.Y-i) ||
			!bishop.coordsOnBoard(bishop.X-i, bishop.Y-i) {
			directions.leftBottom = false
		}
		if bishop.team.FigureExist(bishop.X-i, bishop.Y+i) ||
			bishop.enemy.FigureExist(bishop.X-i, bishop.Y+i) ||
			!bishop.coordsOnBoard(bishop.X-i, bishop.Y+i) {
			directions.leftTop = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (bishop *Bishop) Validation(x int, y int) (bool, string) {
	if !bishop.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if bishop.X == x && bishop.Y == y {
		return false, "can't walk around"
	}
	if bishop.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if bishop.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if is not valid for bishop
	if (x < bishop.X && y < bishop.Y && bishop.X-x != bishop.Y-y) ||
		(x < bishop.X && y > bishop.Y && bishop.X-x != y-bishop.Y) ||
		(x > bishop.X && y < bishop.Y && x-bishop.X != bishop.Y-y) ||
		(x > bishop.X && y > bishop.Y && x-bishop.X != y-bishop.Y) {
		return false, "bishop doesn't walk like that"
	}
	// detect Positions for move and check it for input data move coords
	for _, position := range bishop.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (bishop *Bishop) Move(x int, y int) {
	bishop.team.pawnDoubleMove.clearPawnDoubleMove()
	bishop.MoveFigure(x, y)
}
