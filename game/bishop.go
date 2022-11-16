package game

func NewBishop(x, y int, t *Team) *Bishop {
	f := &Bishop{}
	f.SetName("bishop")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// Bishop is data type of chess figure
type Bishop struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (b *Bishop) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range b.detectionOfBrokenFields() {
		if !b.team.FigureExist(position.X, position.Y) && !b.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (b *Bishop) detectionOfBrokenFields() []Position {
	var data []Position
	directions := struct {
		rightTop    bool
		rightBottom bool
		leftBottom  bool
		leftTop     bool
	}{true, true, true, true}
	for i := 1; i <= 7; i++ {
		if directions.rightTop && b.coordsOnBoard(b.X+i, b.Y+i) {
			data = append(data, Position{X: b.X + i, Y: b.Y + i})
		}
		if directions.rightBottom && b.coordsOnBoard(b.X+i, b.Y-i) {
			data = append(data, Position{X: b.X + i, Y: b.Y - i})
		}
		if directions.leftBottom && b.coordsOnBoard(b.X-i, b.Y-i) {
			data = append(data, Position{X: b.X - i, Y: b.Y - i})
		}
		if directions.leftTop && b.coordsOnBoard(b.X-i, b.Y+i) {
			data = append(data, Position{X: b.X - i, Y: b.Y + i})
		}
		if b.team.FigureExist(b.X+i, b.Y+i) ||
			b.team.enemy.FigureExist(b.X+i, b.Y+i) ||
			!b.coordsOnBoard(b.X+i, b.Y+i) {
			directions.rightTop = false
		}
		if b.team.FigureExist(b.X+i, b.Y-i) ||
			b.team.enemy.FigureExist(b.X+i, b.Y-i) ||
			!b.coordsOnBoard(b.X+i, b.Y-i) {
			directions.rightBottom = false
		}
		if b.team.FigureExist(b.X-i, b.Y-i) ||
			b.team.enemy.FigureExist(b.X-i, b.Y-i) ||
			!b.coordsOnBoard(b.X-i, b.Y-i) {
			directions.leftBottom = false
		}
		if b.team.FigureExist(b.X-i, b.Y+i) ||
			b.team.enemy.FigureExist(b.X-i, b.Y+i) ||
			!b.coordsOnBoard(b.X-i, b.Y+i) {
			directions.leftTop = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (b *Bishop) Validation(x int, y int) (bool, string) {
	if !b.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if b.X == x && b.Y == y {
		return false, "can't walk around"
	}
	if b.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if b.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// if is not valid for bishop
	if (x < b.X && y < b.Y && b.X-x != b.Y-y) ||
		(x < b.X && y > b.Y && b.X-x != y-b.Y) ||
		(x > b.X && y < b.Y && x-b.X != b.Y-y) ||
		(x > b.X && y > b.Y && x-b.X != y-b.Y) {
		return false, "bishop doesn't walk like that"
	}
	// detect Positions for move and check it for input data move coords
	for _, position := range b.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (b *Bishop) Move(x int, y int) {
	b.team.pawnDoubleMove.clearPawnDoubleMove()
	b.MoveFigure(x, y)
}
