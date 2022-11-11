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
func (pDM *pawnDoubleMove) isTakeOnThePass(x, y int) bool {
	if pDM.pawn != nil && pDM.Position.X == x && pDM.Position.Y == y {
		return true
	}
	return false
}

// pawnTakeOnThePass makes a pawn take on the pass
func (pDM *pawnDoubleMove) pawnTakeOnThePass(x, y int) {
	if pDM.Position.X == x && pDM.Position.Y == y {
		// eat figure
		figureID, err := pDM.pawn.team.GetFigureID(pDM.pawn.Position.X, pDM.pawn.Position.Y)
		if err != nil {
			log.Println(err)
		}
		pDM.pawn.team.Eaten[figureID] = pDM.pawn.team.Figures[figureID]
		delete(pDM.pawn.team.Figures, figureID)
	}
}

// clearPawnDoubleMove clear data about double pawn move
func (pDM *pawnDoubleMove) clearPawnDoubleMove() {
	pDM.Position = Position{}
	pDM.pawn = nil
}

// pawnMakesDoubleMove remember data about double pawn move
func (pDM *pawnDoubleMove) pawnMakesDoubleMove(pawn *Pawn, from, to *Position) {
	if to.Y == from.Y+2 || to.Y == from.Y-2 {
		pDM.pawn = pawn
		pDM.Position.X = to.X
		pDM.Position.Y = (to.Y + from.Y) / 2
	}
}
