package game

func NewQueen(pos *Position, t *Team) *Queen {
	q := &Queen{}
	q.figurer = q
	q.SetName("queen")
	q.Position = pos
	q.SetTeam(t)
	return q
}

// Queen is data type of chess figure
type Queen struct {
	Figure
}

// GetBrokenFields return a slice of Positions with broken fields
func (q *Queen) GetBrokenFields() *Positions {
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
	return q.GetPositionsByDirectionsAndMaxRemote(opened, 7)
}

// CanWalkLikeThat desc
func (q *Queen) CanWalkLikeThat(pos *Position) bool {
	if (q.X == pos.X || q.Y == pos.Y) ||
		(pos.X < q.X && pos.Y < q.Y && q.X-pos.X == q.Y-pos.Y) ||
		(pos.X < q.X && pos.Y > q.Y && q.X-pos.X == pos.Y-q.Y) ||
		(pos.X > q.X && pos.Y < q.Y && pos.X-q.X == q.Y-pos.Y) ||
		(pos.X > q.X && pos.Y > q.Y && pos.X-q.X == pos.Y-q.Y) {
		return true
	}
	return false
}
