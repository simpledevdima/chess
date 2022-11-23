package game

// NewKing returns a reference to the new king
// with references to the position and command passed in the argument
func NewKing(pos *Position, t *Team) *King {
	k := &King{}
	k.figurer = k
	k.SetName("king")
	k.Position = pos
	k.SetTeam(t)
	return k
}

// King is data type of chess figure
type King struct {
	Figure
}

// GetBrokenFields return a slice of Positions with broken fields
func (k *King) GetBrokenFields() *Positions {
	opened := map[Direction]bool{
		top:         true,
		topRight:    true,
		right:       true,
		rightBottom: true,
		bottom:      true,
		bottomLeft:  true,
		left:        true,
		leftTop:     true,
	}
	return k.GetPositionsByDirectionsAndMaxRemote(opened, 1)
}

// GetPossibleMoves return slice of Position with coords for possible moves
// has is a boolean variable passed as an argument
// if set to true, returns the map with the first value found, interrupting further calculations
// created in order to minimize the load in case you need to know that there are available moves
func (k *King) GetPossibleMoves(has bool) *Moves {
	mvs := make(Moves)
	var mi MoveIndex

	// castling
	if !k.IsAlreadyMove() {
		pos := NewPosition(3, k.Y)
		if k.castlingIsPossible(pos) && !k.kingOnTheBeatenFieldAfterMove(pos) {
			mi = mvs.Set(mi, NewMove(pos))
		}
		pos = NewPosition(7, k.Y)
		if k.castlingIsPossible(pos) && !k.kingOnTheBeatenFieldAfterMove(pos) {
			mi = mvs.Set(mi, NewMove(pos))
		}
	}

	for _, position := range *k.GetBrokenFields() {
		if !k.team.Figures.ExistsByPosition(position) && !k.kingOnTheBeatenFieldAfterMove(position) {
			mi = mvs.Set(mi, NewMove(position))
			if has {
				return &mvs
			}
		}
	}
	return &mvs
}

// CanWalkLikeThat returns true if the king's move matches the rules for how he moves, otherwise returns false
// this method does not check if the king hit the beaten field after it has been committed
func (k *King) CanWalkLikeThat(pos *Position) bool {
	if (k.X-1 == pos.X || k.X == pos.X || k.X+1 == pos.X) &&
		(k.Y-1 == pos.Y || k.Y == pos.Y || k.Y+1 == pos.Y) {
		return true
	}
	if !k.IsAlreadyMove() && pos.X == 3 || pos.X == 7 {
		return k.castlingIsPossible(pos)
	}
	return false
}

// castlingIsPossible returns true if the castling move matches the rules otherwise returns false
// this method does not check if the king hit the beaten field after it has been committed
func (k *King) castlingIsPossible(pos *Position) bool {
	switch {
	case k.IsAlreadyMove() || k.team.CheckingCheck():
		return false
	case pos.X == 3 || pos.X == 7:
		validation := func(from, to, rookX uint8) bool {
			for x := from; x <= to; x++ {
				if k.team.Figures.ExistsByPosition(NewPosition(x, k.Y)) ||
					k.team.enemy.Figures.ExistsByPosition(NewPosition(x, k.Y)) {
					return false
				}
			}
			rookPos := NewPosition(rookX, k.Y)
			if !k.team.Figures.ExistsByPosition(rookPos) ||
				k.team.Figures.GetByPosition(rookPos).IsAlreadyMove() {
				return false
			}
			return true
		}
		switch pos.X {
		case 3: // long
			return validation(2, 4, 1)
		case 7: // short
			return validation(6, 7, 8)
		}
		fallthrough
	default:
		return false
	}
}
