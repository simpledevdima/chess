package game

// pawnDoubleMove data type containing data and methods for capturing a pawn on the pass
type pawnDoubleMove struct {
	*Position
	pawn *Figure
}

// isTakeOnThePass returns true if it is possible to capture on the pass otherwise returns false
func (p *pawnDoubleMove) isTakeOnThePass(pos *Position) bool {
	if p.pawn != nil && *p.Position == *pos {
		return true
	}
	return false
}

// pawnTakeOnThePass makes a pawn take on the pass
func (p *pawnDoubleMove) pawnTakeOnThePass(pos *Position) {
	if p.Position != nil && *p.Position == *pos {
		// eat figure
		figureID, figure := p.pawn.team.Figures.GetIndexAndFigureByPosition(p.pawn.GetPosition())
		p.pawn.team.Eaten.Set(figureID, figure)
		p.pawn.team.Figures.Remove(figureID)
	}
}

// clearPawnDoubleMove clear data about double pawn move
func (p *pawnDoubleMove) clearPawnDoubleMove() {
	p.Position = nil
	p.pawn = nil
}

// pawnMakesDoubleMove remember data about double pawn move
func (p *pawnDoubleMove) pawnMakesDoubleMove(pawn *Figure, from, to *Position) {
	p.clearPawnDoubleMove()
	if to.Y == from.Y+2 || to.Y == from.Y-2 {
		p.pawn = pawn
		p.Position = NewPosition(to.X, (to.Y+from.Y)/2)
	}
}
