package game

import "log"

func NewPawn(x, y int, t *Team) *Pawn {
	f := &Pawn{}
	f.SetName("pawn")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// Pawn is data type of chess figure
type Pawn struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (p *Pawn) DetectionOfPossibleMove() []*Position {
	var data []*Position
	switch p.team.Name {
	case White:
		if p.coordsOnBoard(p.X, p.Y+1) &&
			!p.kingOnTheBeatenFieldAfterMove(p.X, p.Y+1) &&
			!p.team.FigureExist(p.X, p.Y+1) &&
			!p.team.enemy.FigureExist(p.X, p.Y+1) {
			data = append(data, NewPosition(p.X, p.Y+1))
		}
		if p.coordsOnBoard(p.X, p.Y+2) &&
			!p.kingOnTheBeatenFieldAfterMove(p.X, p.Y+2) &&
			!p.IsAlreadyMove() &&
			!p.team.FigureExist(p.X, p.Y+1) &&
			!p.team.FigureExist(p.X, p.Y+2) &&
			!p.team.enemy.FigureExist(p.X, p.Y+1) &&
			!p.team.enemy.FigureExist(p.X, p.Y+2) {
			data = append(data, NewPosition(p.X, p.Y+2))
		}
	case Black:
		if p.coordsOnBoard(p.X, p.Y-1) &&
			!p.kingOnTheBeatenFieldAfterMove(p.X, p.Y-1) &&
			!p.team.FigureExist(p.X, p.Y-1) &&
			!p.team.enemy.FigureExist(p.X, p.Y-1) {
			data = append(data, NewPosition(p.X, p.Y-1))
		}
		if p.coordsOnBoard(p.X, p.Y-2) &&
			!p.kingOnTheBeatenFieldAfterMove(p.X, p.Y-2) &&
			!p.IsAlreadyMove() &&
			!p.team.FigureExist(p.X, p.Y-1) &&
			!p.team.FigureExist(p.X, p.Y-2) &&
			!p.team.enemy.FigureExist(p.X, p.Y-1) &&
			!p.team.enemy.FigureExist(p.X, p.Y-2) {
			data = append(data, NewPosition(p.X, p.Y-2))
		}
	}
	for _, position := range p.detectionOfBrokenFields() {
		if (p.team.enemy.FigureExist(position.X, position.Y) ||
			p.team.enemy.pawnDoubleMove.isTakeOnThePass(position.X, position.Y)) &&
			!p.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			data = append(data, position)
		}
	}
	return data
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (p *Pawn) detectionOfBrokenFields() []*Position {
	var data []*Position
	switch p.team.Name {
	case White:
		if p.coordsOnBoard(p.X+1, p.Y+1) {
			data = append(data, NewPosition(p.X+1, p.Y+1))
		}
		if p.coordsOnBoard(p.X-1, p.Y+1) {
			data = append(data, NewPosition(p.X-1, p.Y+1))
		}
	case Black:
		if p.coordsOnBoard(p.X+1, p.Y-1) {
			data = append(data, NewPosition(p.X+1, p.Y-1))
		}
		if p.coordsOnBoard(p.X-1, p.Y-1) {
			data = append(data, NewPosition(p.X-1, p.Y-1))
		}
	}
	return data
}

// Validation return true if this move are valid or return false
func (p *Pawn) Validation(x int, y int) (bool, string) {
	if !p.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if p.X == x && p.Y == y {
		return false, "can't walk around"
	}
	if p.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if p.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for eat and check it for input data eat coords
	for _, position := range p.detectionOfBrokenFields() {
		if position.X == x && position.Y == y &&
			(p.team.enemy.FigureExist(x, y) || p.team.enemy.pawnDoubleMove.isTakeOnThePass(x, y)) {
			return true, ""
		}
	}
	// move pawn
	for _, Position := range p.DetectionOfPossibleMove() {
		if Position.X == x && Position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (p *Pawn) Move(x int, y int) {
	p.team.enemy.pawnDoubleMove.pawnTakeOnThePass(x, y)
	p.team.pawnDoubleMove.clearPawnDoubleMove()
	p.team.pawnDoubleMove.pawnMakesDoubleMove(p, &Position{X: p.X, Y: p.Y}, &Position{X: x, Y: y})
	p.MoveFigure(x, y)
	p.transformPawnTOQueen(x, y)
}

// transformPawnToQueen promote a pawn to a queen
func (p *Pawn) transformPawnTOQueen(x, y int) {
	if p.Y == 1 || p.Y == 8 {
		figureID, err := p.team.GetFigureID(x, y)
		if err != nil {
			log.Println(err)
		}
		// replace pawn to queen
		p.team.Figures[figureID] = NewQueen(x, y, p.team)
		p.team.Figures[figureID].setAlreadyMove(true)
	}
}
