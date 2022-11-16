package game

func NewQueen(x, y int, t *Team) *Queen {
	f := &Queen{}
	f.SetName("queen")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// Queen is data type of chess figure
type Queen struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (q *Queen) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range q.detectionOfBrokenFields() {
		if !q.team.FigureExist(position.X, position.Y) && !q.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (q *Queen) detectionOfBrokenFields() []Position {
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
		if directions.top && q.coordsOnBoard(q.X, q.Y+i) {
			data = append(data, Position{X: q.X, Y: q.Y + i})
		}
		if directions.rightTop && q.coordsOnBoard(q.X+i, q.Y+i) {
			data = append(data, Position{X: q.X + i, Y: q.Y + i})
		}
		if directions.right && q.coordsOnBoard(q.X+i, q.Y) {
			data = append(data, Position{X: q.X + i, Y: q.Y})
		}
		if directions.rightBottom && q.coordsOnBoard(q.X+i, q.Y-i) {
			data = append(data, Position{X: q.X + i, Y: q.Y - i})
		}
		if directions.bottom && q.coordsOnBoard(q.X, q.Y-i) {
			data = append(data, Position{X: q.X, Y: q.Y - i})
		}
		if directions.leftBottom && q.coordsOnBoard(q.X-i, q.Y-i) {
			data = append(data, Position{X: q.X - i, Y: q.Y - i})
		}
		if directions.left && q.coordsOnBoard(q.X-i, q.Y) {
			data = append(data, Position{X: q.X - i, Y: q.Y})
		}
		if directions.leftTop && q.coordsOnBoard(q.X-i, q.Y+i) {
			data = append(data, Position{X: q.X - i, Y: q.Y + i})
		}
		if q.team.FigureExist(q.X, q.Y+i) ||
			q.team.enemy.FigureExist(q.X, q.Y+i) ||
			!q.coordsOnBoard(q.X, q.Y+i) {
			directions.top = false
		}
		if q.team.FigureExist(q.X+i, q.Y+i) ||
			q.team.enemy.FigureExist(q.X+i, q.Y+i) ||
			!q.coordsOnBoard(q.X+i, q.Y+i) {
			directions.rightTop = false
		}
		if q.team.FigureExist(q.X+i, q.Y) ||
			q.team.enemy.FigureExist(q.X+i, q.Y) ||
			!q.coordsOnBoard(q.X+i, q.Y) {
			directions.right = false
		}
		if q.team.FigureExist(q.X+i, q.Y-i) ||
			q.team.enemy.FigureExist(q.X+i, q.Y-i) ||
			!q.coordsOnBoard(q.X+i, q.Y-i) {
			directions.rightBottom = false
		}
		if q.team.FigureExist(q.X, q.Y-i) ||
			q.team.enemy.FigureExist(q.X, q.Y-i) ||
			!q.coordsOnBoard(q.X, q.Y-i) {
			directions.bottom = false
		}
		if q.team.FigureExist(q.X-i, q.Y-i) ||
			q.team.enemy.FigureExist(q.X-i, q.Y-i) ||
			!q.coordsOnBoard(q.X-i, q.Y-i) {
			directions.leftBottom = false
		}
		if q.team.FigureExist(q.X-i, q.Y) ||
			q.team.enemy.FigureExist(q.X-i, q.Y) ||
			!q.coordsOnBoard(q.X-i, q.Y) {
			directions.left = false
		}
		if q.team.FigureExist(q.X-i, q.Y+i) ||
			q.team.enemy.FigureExist(q.X-i, q.Y+i) ||
			!q.coordsOnBoard(q.X-i, q.Y+i) {
			directions.leftTop = false
		}
	}

	return data
}

// Validation return true if this move are valid or return false
func (q *Queen) Validation(x int, y int) (bool, string) {
	if !q.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if q.X == x && q.Y == y {
		return false, "can't walk around"
	}
	if q.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if q.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if change x && y is not valid for queen
	if (q.X != x && q.Y != y && x < q.X && y < q.Y && q.X-x != q.Y-y) ||
		(q.X != x && q.Y != y && x < q.X && y > q.Y && q.X-x != y-q.Y) ||
		(q.X != x && q.Y != y && x > q.X && y < q.Y && x-q.X != q.Y-y) ||
		(q.X != x && q.Y != y && x > q.X && y > q.Y && x-q.X != y-q.Y) {
		return false, "queen doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range q.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (q *Queen) Move(x int, y int) {
	q.team.pawnDoubleMove.clearPawnDoubleMove()
	q.MoveFigure(x, y)
}
