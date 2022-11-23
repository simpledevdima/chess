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
func (k *King) GetPossibleMoves(has bool) *Positions {
	poss := make(Positions)
	var pi PositionIndex

	// add castling to possible moves if those moves are possible
	if !k.IsAlreadyMove() && !k.team.CheckingCheck() {
		func() {
			// in the long side
			for x := uint8(1); x <= 3; x++ {
				if k.team.Figures.ExistsByPosition(NewPosition(k.X-x, k.Y)) ||
					k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X-x, k.Y)) {
					return
				}
			}
			pos := NewPosition(3, k.Y)
			if !k.team.Figures.ExistsByPosition(NewPosition(k.X-4, k.Y)) ||
				k.team.Figures.GetByPosition(NewPosition(k.X-4, k.Y)).IsAlreadyMove() ||
				k.kingOnTheBeatenFieldAfterMove(pos) {
				return
			}
			pi = poss.Set(pi, pos)
		}()
		if has && len(poss) > 0 {
			return &poss
		}
		func() {
			// in the short side
			for x := uint8(1); x <= 2; x++ {
				if k.team.Figures.ExistsByPosition(NewPosition(k.X+x, k.Y)) ||
					k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X+x, k.Y)) {
					return
				}
			}
			pos := NewPosition(7, k.Y)
			if !k.team.Figures.ExistsByPosition(NewPosition(k.X+3, k.Y)) ||
				k.team.Figures.GetByPosition(NewPosition(k.X+3, k.Y)).IsAlreadyMove() ||
				k.kingOnTheBeatenFieldAfterMove(pos) {
				return
			}
			pi = poss.Set(pi, pos)
		}()
		if has && len(poss) > 0 {
			return &poss
		}
	}

	for _, position := range *k.GetBrokenFields() {
		if !k.team.Figures.ExistsByPosition(position) && !k.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
			if has {
				return &poss
			}
		}
	}
	return &poss
}

// CanWalkLikeThat returns true if the king's move matches the rules for how he moves, otherwise returns false
func (k *King) CanWalkLikeThat(pos *Position) bool {
	if (k.X-1 == pos.X || k.X == pos.X || k.X+1 == pos.X) &&
		(k.Y-1 == pos.Y || k.Y == pos.Y || k.Y+1 == pos.Y) {
		return true
	}
	return false
}
