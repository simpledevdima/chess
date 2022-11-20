package game

func NewKing(pos *Position, t *Team) *King {
	k := &King{}
	k.SetName("king")
	k.Position = pos
	k.SetTeam(t)
	return k
}

// King is data type of chess figure
type King struct {
	Figure
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (k *King) GetPossibleMoves() *Positions {
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
func (k *King) GetBrokenFields() *Positions {
	poss := make(Positions)
	var pi PositionIndex

	pos := NewPosition(k.X, k.Y+1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+1, k.Y+1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+1, k.Y)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X+1, k.Y-1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X, k.Y-1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-1, k.Y-1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-1, k.Y)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	pos = NewPosition(k.X-1, k.Y+1)
	if k.positionOnBoard(pos) {
		pi = poss.Set(pi, pos)
	}

	return &poss
}

// Validation return true if this move are valid or return false
func (k *King) Validation(pos *Position) (bool, string) {
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
	// castling
	if !k.alreadyMove {
		if pos.X == 3 {
			if !k.team.CheckingCheck() &&
				!k.team.Figures.ExistsByPosition(NewPosition(k.X-1, k.Y)) && !k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X-1, k.Y)) &&
				!k.team.Figures.ExistsByPosition(NewPosition(k.X-2, k.Y)) && !k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X-2, k.Y)) &&
				!k.team.Figures.ExistsByPosition(NewPosition(k.X-3, k.Y)) && !k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X-3, k.Y)) &&
				k.team.Figures.ExistsByPosition(NewPosition(k.X-4, k.Y)) {
				if !k.team.Figures.GetByPosition(NewPosition(k.X-4, k.Y)).IsAlreadyMove() {
					return true, ""
				}
			}
		} else if pos.X == 7 {
			if !k.team.CheckingCheck() &&
				!k.team.Figures.ExistsByPosition(NewPosition(k.X+1, k.Y)) && !k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X+1, k.Y)) &&
				!k.team.Figures.ExistsByPosition(NewPosition(k.X+2, k.Y)) && !k.team.enemy.Figures.ExistsByPosition(NewPosition(k.X+2, k.Y)) &&
				k.team.Figures.ExistsByPosition(NewPosition(k.X+3, k.Y)) {
				if !k.team.Figures.GetByPosition(NewPosition(k.X+3, k.Y)).IsAlreadyMove() {
					return true, ""
				}
			}
		}
	}
	// detect Position for move and check it for input data move coords
	for _, position := range *k.GetBrokenFields() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}
