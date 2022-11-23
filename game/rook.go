package game

// NewRook returns a reference to the new rook
// with references to the position and command passed in the argument
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

// CanWalkLikeThat returns true if the rook's move follows the rules for how it moves, otherwise it returns false
// this method does not check if the king hit the beaten field after it has been committed
func (r *Rook) CanWalkLikeThat(pos *Position) bool {
	if r.X == pos.X || r.Y == pos.Y {
		return true
	}
	return false
}
