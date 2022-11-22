package game

func NewKnight(pos *Position, t *Team) *Knight {
	k := &Knight{}
	k.figurer = k
	k.SetName("knight")
	k.Position = pos
	k.SetTeam(t)
	return k
}

// Knight is data type of chess figure
type Knight struct {
	Figure
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

// CanWalkLikeThat desc
func (k *Knight) CanWalkLikeThat(pos *Position) bool {
	if (k.X+1 == pos.X && k.Y+2 == pos.Y) ||
		(k.X+2 == pos.X && k.Y+1 == pos.Y) ||
		(k.X-1 == pos.X && k.Y-2 == pos.Y) ||
		(k.X-2 == pos.X && k.Y-1 == pos.Y) ||
		(k.X+1 == pos.X && k.Y-2 == pos.Y) ||
		(k.X+2 == pos.X && k.Y-1 == pos.Y) ||
		(k.X-1 == pos.X && k.Y+2 == pos.Y) ||
		(k.X-2 == pos.X && k.Y+1 == pos.Y) {
		return true
	}
	return false
}
