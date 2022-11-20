package game

func NewBishop(pos *Position, t *Team) *Bishop {
	b := &Bishop{}
	b.Position = pos
	b.SetName("bishop")
	b.SetTeam(t)
	return b
}

// Bishop is data type of chess figure
type Bishop struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (b *Bishop) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range b.detectionOfBrokenFields() {
		if !b.team.Figures.ExistsByPosition(position) && !b.kingOnTheBeatenFieldAfterMove(position) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (b *Bishop) detectionOfBrokenFields() []*Position {
	var data []*Position
	directions := struct {
		rightTop    bool
		rightBottom bool
		leftBottom  bool
		leftTop     bool
	}{true, true, true, true}
	for i := 1; i <= 7; i++ {
		pos := NewPosition(b.X+i, b.Y+i)
		if directions.rightTop && b.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if b.team.Figures.ExistsByPosition(pos) ||
			b.team.enemy.Figures.ExistsByPosition(pos) ||
			!b.positionOnBoard(pos) {
			directions.rightTop = false
		}

		pos = NewPosition(b.X+i, b.Y-i)
		if directions.rightBottom && b.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if b.team.Figures.ExistsByPosition(pos) ||
			b.team.enemy.Figures.ExistsByPosition(pos) ||
			!b.positionOnBoard(pos) {
			directions.rightBottom = false
		}

		pos = NewPosition(b.X-i, b.Y-i)
		if directions.leftBottom && b.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if b.team.Figures.ExistsByPosition(pos) ||
			b.team.enemy.Figures.ExistsByPosition(pos) ||
			!b.positionOnBoard(pos) {
			directions.leftBottom = false
		}

		pos = NewPosition(b.X-i, b.Y+i)
		if directions.leftTop && b.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if b.team.Figures.ExistsByPosition(pos) ||
			b.team.enemy.Figures.ExistsByPosition(pos) ||
			!b.positionOnBoard(pos) {
			directions.leftTop = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (b *Bishop) Validation(pos *Position) (bool, string) {
	if !b.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *b.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if b.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if b.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	// if is not valid for bishop
	x, y := pos.Get()
	if (x < b.X && y < b.Y && b.X-x != b.Y-y) ||
		(x < b.X && y > b.Y && b.X-x != y-b.Y) ||
		(x > b.X && y < b.Y && x-b.X != b.Y-y) ||
		(x > b.X && y > b.Y && x-b.X != y-b.Y) {
		return false, "bishop doesn't walk like that"
	}
	// detect Positions for move and check it for input data move coords
	for _, position := range b.detectionOfBrokenFields() {
		if *position == *pos {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (b *Bishop) Move(pos *Position) {
	b.team.pawnDoubleMove.clearPawnDoubleMove()
	b.MoveFigure(pos)
}
