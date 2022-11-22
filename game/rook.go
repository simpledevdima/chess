package game

func NewRook(pos *Position, t *Team) *Rook {
	r := &Rook{}
	r.figurer = r
	r.SetName("rook")
	r.Position = pos
	r.SetTeam(t)
	return r
}

// Rook is data type of chess figure
type Rook struct {
	Figure
}

// GetBrokenFields return a slice of Positions with broken fields
func (r *Rook) GetBrokenFields() *Positions {
	opened := map[Direction]bool{
		top:    true,
		right:  true,
		bottom: true,
		left:   true,
	}
	return r.GetPositionsByDirectionsAndMaxRemote(opened, 7)
}

// CanWalkLikeThat desc
func (r *Rook) CanWalkLikeThat(pos *Position) bool {
	if r.X == pos.X || r.Y == pos.Y {
		return true
	}
	return false
}
