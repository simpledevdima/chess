package game

func NewBishop(pos *Position, t *Team) *Bishop {
	b := &Bishop{}
	b.figurer = b
	b.Position = pos
	b.SetName("bishop")
	b.SetTeam(t)
	return b
}

// Bishop is data type of chess figure
type Bishop struct {
	Figure
}

// GetBrokenFields return a slice of Positions with broken fields
func (b *Bishop) GetBrokenFields() *Positions {
	opened := map[Direction]bool{
		topRight:    true,
		rightBottom: true,
		bottomLeft:  true,
		leftTop:     true,
	}
	return b.GetPositionsByDirectionsAndMaxRemote(opened, 7)
}

// CanWalkLikeThat desc
func (b *Bishop) CanWalkLikeThat(pos *Position) bool {
	if (pos.X < b.X && pos.Y < b.Y && b.X-pos.X == b.Y-pos.Y) ||
		(pos.X < b.X && pos.Y > b.Y && b.X-pos.X == pos.Y-b.Y) ||
		(pos.X > b.X && pos.Y < b.Y && pos.X-b.X == b.Y-pos.Y) ||
		(pos.X > b.X && pos.Y > b.Y && pos.X-b.X == pos.Y-b.Y) {
		return true
	}
	return false
}
