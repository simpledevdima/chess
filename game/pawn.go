package game

import (
	"log"
)

// Pawn is data type of chess figure
type Pawn struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (pawn *Pawn) DetectionOfPossibleMove() []Position {
	var data []Position
	switch pawn.team.Name {
	case White:
		if pawn.coordsOnBoard(pawn.X, pawn.Y+1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y+1) {
			data = append(data, Position{X: pawn.X, Y: pawn.Y + 1})
		}
		if pawn.coordsOnBoard(pawn.X, pawn.Y+2) &&
			!pawn.IsAlreadyMove() &&
			!pawn.team.FigureExist(pawn.X, pawn.Y+1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y+1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y+2) {
			data = append(data, Position{X: pawn.X, Y: pawn.Y + 2})
		}
	case Black:
		if pawn.coordsOnBoard(pawn.X, pawn.Y-1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y-1) {
			data = append(data, Position{X: pawn.X, Y: pawn.Y - 1})
		}
		if pawn.coordsOnBoard(pawn.X, pawn.Y-2) &&
			!pawn.IsAlreadyMove() &&
			!pawn.team.FigureExist(pawn.X, pawn.Y-1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y-1) &&
			!pawn.enemy.FigureExist(pawn.X, pawn.Y-2) {
			data = append(data, Position{X: pawn.X, Y: pawn.Y - 2})
		}
	}
	return data
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (pawn *Pawn) detectionOfBrokenFields() []Position {
	var data []Position
	switch pawn.team.Name {
	case White:
		if pawn.coordsOnBoard(pawn.X+1, pawn.Y+1) {
			data = append(data, Position{X: pawn.X + 1, Y: pawn.Y + 1})
		}
		if pawn.coordsOnBoard(pawn.X-1, pawn.Y+1) {
			data = append(data, Position{X: pawn.X - 1, Y: pawn.Y + 1})
		}
	case Black:
		if pawn.coordsOnBoard(pawn.X+1, pawn.Y-1) {
			data = append(data, Position{X: pawn.X + 1, Y: pawn.Y - 1})
		}
		if pawn.coordsOnBoard(pawn.X-1, pawn.Y-1) {
			data = append(data, Position{X: pawn.X - 1, Y: pawn.Y - 1})
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (pawn *Pawn) Validation(x int, y int) (bool, string) {
	if !pawn.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if pawn.X == x && pawn.Y == y {
		return false, "can't walk around"
	}
	if pawn.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if pawn.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for eat and check it for input data eat coords
	for _, Position := range pawn.detectionOfBrokenFields() {
		if Position.X == x && Position.Y == y &&
			(pawn.enemy.FigureExist(x, y) || pawn.team.pawnDoubleMove.isTakeOnThePass(x, y)) {
			return true, ""
		}
	}
	// move pawn
	for _, Position := range pawn.DetectionOfPossibleMove() {
		if Position.X == x && Position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (pawn *Pawn) Move(x int, y int) {
	pawn.team.pawnDoubleMove.pawnTakeOnThePass(x, y)
	pawn.team.pawnDoubleMove.clearPawnDoubleMove()
	pawn.team.pawnDoubleMove.pawnMakesDoubleMove(pawn, &Position{X: pawn.X, Y: pawn.Y}, &Position{X: x, Y: y})

	pawn.MoveFigure(x, y)

	// transform pawn to queen
	if pawn.Y == 1 || pawn.Y == 8 {
		figureID, err := pawn.team.GetFigureID(x, y)
		if err != nil {
			log.Println(err)
		}
		// replace pawn to queen
		pawn.team.Figures[figureID] = &Queen{}
		pawn.team.Figures[figureID].SetPosition(x, y)
		pawn.team.Figures[figureID].SetTeams(pawn.team, pawn.enemy)
		pawn.team.Figures[figureID].setAlreadyMove(true)
	}
}
