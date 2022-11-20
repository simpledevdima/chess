package game

func NewQueen(pos *Position, t *Team) *Queen {
	q := &Queen{}
	q.SetName("queen")
	q.Position = pos
	q.SetTeam(t)
	return q
}

// Queen is data type of chess figure
type Queen struct {
	Figure
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (q *Queen) GetPossibleMoves() *Positions {
	poss := make(Positions)
	var pi PositionIndex
	for _, position := range *q.GetBrokenFields() {
		if !q.team.Figures.ExistsByPosition(position) && !q.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
		}
	}
	return &poss
}

// GetBrokenFields return a slice of Positions with broken fields
func (q *Queen) GetBrokenFields() *Positions {
	poss := make(Positions)
	var pi PositionIndex

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
	for i := uint8(1); i <= 7; i++ {
		pos := NewPosition(q.X, q.Y+i)
		if directions.top && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.top = false
		}

		pos = NewPosition(q.X+i, q.Y+i)
		if directions.rightTop && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.rightTop = false
		}

		pos = NewPosition(q.X+i, q.Y)
		if directions.right && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.right = false
		}

		pos = NewPosition(q.X+i, q.Y-i)
		if directions.rightBottom && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.rightBottom = false
		}

		pos = NewPosition(q.X, q.Y-i)
		if directions.bottom && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.bottom = false
		}

		pos = NewPosition(q.X-i, q.Y-i)
		if directions.leftBottom && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.leftBottom = false
		}

		pos = NewPosition(q.X-i, q.Y)
		if directions.left && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.left = false
		}

		pos = NewPosition(q.X-i, q.Y+i)
		if directions.leftTop && q.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		if q.team.Figures.ExistsByPosition(pos) ||
			q.team.enemy.Figures.ExistsByPosition(pos) ||
			!q.positionOnBoard(pos) {
			directions.leftTop = false
		}
	}

	return &poss
}

// Validation return true if this move are valid or return false
func (q *Queen) Validation(pos *Position) (bool, string) {
	if !q.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *q.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if q.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if q.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	x, y := pos.Get()
	// if change x && y is not valid for queen
	if (q.X != x && q.Y != y && x < q.X && y < q.Y && q.X-x != q.Y-y) ||
		(q.X != x && q.Y != y && x < q.X && y > q.Y && q.X-x != y-q.Y) ||
		(q.X != x && q.Y != y && x > q.X && y < q.Y && x-q.X != q.Y-y) ||
		(q.X != x && q.Y != y && x > q.X && y > q.Y && x-q.X != y-q.Y) {
		return false, "queen doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range *q.GetBrokenFields() {
		if position.X == x && position.Y == y {
			// this move is valid
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}
