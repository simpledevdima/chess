package game

func NewRook(pos *Position, t *Team) *Rook {
	r := &Rook{}
	r.SetName("rook")
	r.Position = pos
	r.SetTeam(t)
	return r
}

// Rook is data type of chess figure
type Rook struct {
	Figure
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (r *Rook) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range r.detectionOfBrokenFields() {
		if !r.team.Figures.ExistsByPosition(position) && !r.kingOnTheBeatenFieldAfterMove(position) {
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
	for i := uint8(1); i <= 7; i++ {
		pos := NewPosition(r.X, r.Y+i)
		if directions.top && r.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if r.team.Figures.ExistsByPosition(pos) ||
			r.team.enemy.Figures.ExistsByPosition(pos) ||
			!r.positionOnBoard(pos) {
			directions.top = false
		}

		pos = NewPosition(r.X+i, r.Y)
		if directions.right && r.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if r.team.Figures.ExistsByPosition(pos) ||
			r.team.enemy.Figures.ExistsByPosition(pos) ||
			!r.positionOnBoard(pos) {
			directions.right = false
		}

		pos = NewPosition(r.X, r.Y-i)
		if directions.bottom && r.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if r.team.Figures.ExistsByPosition(pos) ||
			r.team.enemy.Figures.ExistsByPosition(pos) ||
			!r.positionOnBoard(pos) {
			directions.bottom = false
		}

		pos = NewPosition(r.X-i, r.Y)
		if directions.left && r.positionOnBoard(pos) {
			data = append(data, pos)
		}
		if r.team.Figures.ExistsByPosition(pos) ||
			r.team.enemy.Figures.ExistsByPosition(pos) ||
			!r.positionOnBoard(pos) {
			directions.left = false
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (r *Rook) Validation(pos *Position) (bool, string) {
	if !r.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *r.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if r.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if r.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	// if change x && y is not valid for rook
	if r.X != pos.X && r.Y != pos.Y {
		return false, "rook doesn't walk like that"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range r.detectionOfBrokenFields() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}
