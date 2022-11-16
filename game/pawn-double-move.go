package game

import (
	"log"
)

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
		figureID, err := p.pawn.team.GetFigureID(p.pawn.Position.X, p.pawn.Position.Y)
		if err != nil {
			log.Println(err)
		}
		p.pawn.team.Eaten[figureID] = p.pawn.team.Figures[figureID]
		delete(p.pawn.team.Figures, figureID)
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
