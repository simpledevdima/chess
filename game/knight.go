package game

func NewKnight(pos *Position, t *Team) *Knight {
	k := &Knight{}
	k.SetName("knight")
	k.Position = pos
	k.SetTeam(t)
	return k
}

// Knight is data type of chess figure
type Knight struct {
	Figure
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (k *Knight) GetPossibleMoves() *Positions {
	poss := make(Positions)
	var pi PositionIndex
	for _, position := range *k.GetBrokenFields() {
		if !k.team.Figures.ExistsByPosition(position) && !k.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
		}
	}
	return &poss
}

// GetBrokenFields return a slice of Positions with broken fields
func (k *Knight) GetBrokenFields() *Positions {
	poss := make(Positions)
	var pi PositionIndex

	pos := NewPosition(k.X+1, k.Y+2)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+2, k.Y+1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+2, k.Y-1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+1, k.Y-2)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-1, k.Y-2)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-2, k.Y-1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-2, k.Y+1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-1, k.Y+2)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	return &poss
}

// Validation return true if this move are valid or return false
func (k *Knight) Validation(pos *Position) (bool, string) {
	if !k.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *k.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if k.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if k.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for move and check it for input data move coords
	for _, position := range *k.GetBrokenFields() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}
