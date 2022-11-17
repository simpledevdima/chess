package game

// pawnDoubleMove data type containing data and methods for capturing a pawn on the pass
type pawnDoubleMove struct {
	Position Position
	pawn     *Pawn
}

// isTakeOnThePass returns true if it is possible to capture on the pass otherwise returns false
func (p *pawnDoubleMove) isTakeOnThePass(x, y int) bool {
	if p.pawn != nil && p.Position.X == x && p.Position.Y == y {
		return true
	}
	return false
}

// pawnTakeOnThePass makes a pawn take on the pass
func (p *pawnDoubleMove) pawnTakeOnThePass(x, y int) {
	if p.Position.X == x && p.Position.Y == y {
		// eat figure
		figureID, figure := p.pawn.team.Figures.GetIndexAndFigureByCoords(p.pawn.Position.X, p.pawn.Position.Y)
		p.pawn.team.Eaten.Set(figureID, figure)
		p.pawn.team.Figures.RemoveByIndex(figureID)
	}
}

// clearPawnDoubleMove clear data about double pawn move
func (p *pawnDoubleMove) clearPawnDoubleMove() {
	p.Position = Position{}
	p.pawn = nil
}

// pawnMakesDoubleMove remember data about double pawn move
func (p *pawnDoubleMove) pawnMakesDoubleMove(pawn *Pawn, from, to *Position) {
	if to.Y == from.Y+2 || to.Y == from.Y-2 {
		p.pawn = pawn
		p.Position.X = to.X
		p.Position.Y = (to.Y + from.Y) / 2
	}
}
