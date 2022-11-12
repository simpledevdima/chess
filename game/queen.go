package game

// Queen is data type of chess figure
type Queen struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (queen *Queen) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range queen.detectionOfBrokenFields() {
		if !queen.team.FigureExist(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (queen *Queen) detectionOfBrokenFields() []Position {
	var data []Position

	directions := struct {
		top         bool
		rightTop    bool
		right       bool
		rightBottom bool
		bottom      bool
		leftBottom  bool
		left        bool
		leftTop     bool
	}{true, true, true, true, true, true, true, true}
	for i := 1; i <= 7; i++ {
		if directions.top && queen.coordsOnBoard(queen.X, queen.Y+i) {
			data = append(data, Position{X: queen.X, Y: queen.Y + i})
		}
		if directions.rightTop && queen.coordsOnBoard(queen.X+i, queen.Y+i) {
			data = append(data, Position{X: queen.X + i, Y: queen.Y + i})
		}
		if directions.right && queen.coordsOnBoard(queen.X+i, queen.Y) {
			data = append(data, Position{X: queen.X + i, Y: queen.Y})
		}
		if directions.rightBottom && queen.coordsOnBoard(queen.X+i, queen.Y-i) {
			data = append(data, Position{X: queen.X + i, Y: queen.Y - i})
		}
		if directions.bottom && queen.coordsOnBoard(queen.X, queen.Y-i) {
			data = append(data, Position{X: queen.X, Y: queen.Y - i})
		}
		if directions.leftBottom && queen.coordsOnBoard(queen.X-i, queen.Y-i) {
			data = append(data, Position{X: queen.X - i, Y: queen.Y - i})
		}
		if directions.left && queen.coordsOnBoard(queen.X-i, queen.Y) {
			data = append(data, Position{X: queen.X - i, Y: queen.Y})
		}
		if directions.leftTop && queen.coordsOnBoard(queen.X-i, queen.Y+i) {
			data = append(data, Position{X: queen.X - i, Y: queen.Y + i})
		}
		if queen.team.FigureExist(queen.X, queen.Y+i) ||
			queen.enemy.FigureExist(queen.X, queen.Y+i) ||
			!queen.coordsOnBoard(queen.X, queen.Y+i) {
			directions.top = false
		}
		if queen.team.FigureExist(queen.X+i, queen.Y+i) ||
			queen.enemy.FigureExist(queen.X+i, queen.Y+i) ||
			!queen.coordsOnBoard(queen.X+i, queen.Y+i) {
			directions.rightTop = false
		}
		if queen.team.FigureExist(queen.X+i, queen.Y) ||
			queen.enemy.FigureExist(queen.X+i, queen.Y) ||
			!queen.coordsOnBoard(queen.X+i, queen.Y) {
			directions.right = false
		}
		if queen.team.FigureExist(queen.X+i, queen.Y-i) ||
			queen.enemy.FigureExist(queen.X+i, queen.Y-i) ||
			!queen.coordsOnBoard(queen.X+i, queen.Y-i) {
			directions.rightBottom = false
		}
		if queen.team.FigureExist(queen.X, queen.Y-i) ||
			queen.enemy.FigureExist(queen.X, queen.Y-i) ||
			!queen.coordsOnBoard(queen.X, queen.Y-i) {
			directions.bottom = false
		}
		if queen.team.FigureExist(queen.X-i, queen.Y-i) ||
			queen.enemy.FigureExist(queen.X-i, queen.Y-i) ||
			!queen.coordsOnBoard(queen.X-i, queen.Y-i) {
			directions.leftBottom = false
		}
		if queen.team.FigureExist(queen.X-i, queen.Y) ||
			queen.enemy.FigureExist(queen.X-i, queen.Y) ||
			!queen.coordsOnBoard(queen.X-i, queen.Y) {
			directions.left = false
		}
		if queen.team.FigureExist(queen.X-i, queen.Y+i) ||
			queen.enemy.FigureExist(queen.X-i, queen.Y+i) ||
			!queen.coordsOnBoard(queen.X-i, queen.Y+i) {
			directions.leftTop = false
		}
	}

	return data
}

// Validation return true if this move are valid or return false
func (queen *Queen) Validation(x int, y int) (bool, string) {
	if !queen.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if queen.X == x && queen.Y == y {
		return false, "can't walk around"
	}
	if queen.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if queen.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if change x && y is not valid for queen
	if (queen.X != x && queen.Y != y && x < queen.X && y < queen.Y && queen.X-x != queen.Y-y) ||
		(queen.X != x && queen.Y != y && x < queen.X && y > queen.Y && queen.X-x != y-queen.Y) ||
		(queen.X != x && queen.Y != y && x > queen.X && y < queen.Y && x-queen.X != queen.Y-y) ||
		(queen.X != x && queen.Y != y && x > queen.X && y > queen.Y && x-queen.X != y-queen.Y) {
		return false, "queen doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range queen.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (queen *Queen) Move(x int, y int) {
	queen.team.pawnDoubleMove.clearPawnDoubleMove()
	queen.MoveFigure(x, y)
}
