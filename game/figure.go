package game

import (
	"log"
)

// Figure information about each figure
type Figure struct {
	Name          string `json:"name"`
	*Position     `json:"position"`
	alreadyMove   bool
	team          *Team
	possibleMoves Positions
	brokenFields  Positions
}

// GetName get name from figure
func (f *Figure) GetName() string {
	return f.Name
}

// SetName set name to the figure
func (f *Figure) SetName(name string) {
	f.Name = name
}

// IsAlreadyMove returns true if the figure has moved in the current match otherwise returns false
func (f *Figure) IsAlreadyMove() bool {
	return f.alreadyMove
}

// setAlreadyMove setting a value about moving a figure in the current match
func (f *Figure) setAlreadyMove(flag bool) {
	f.alreadyMove = flag
}

// kingOnTheBeatenFieldAfterMove returns true if your king is on the beaten square after the move otherwise return false
func (f *Figure) kingOnTheBeatenFieldAfterMove(pos *Position) bool {
	curPos := f.Position
	f.SetPosition(pos)
	undoMove := func() {
		f.SetPosition(curPos)
	}
	undoEating := func() {}
	if f.team.enemy.Figures.ExistsByPosition(pos) {
		eatenID := f.team.enemy.Figures.GetIndexByPosition(pos)
		eatenFigure := f.team.enemy.Figures.Get(eatenID)
		f.team.enemy.Figures.Remove(eatenID)
		undoEating = func() {
			f.team.enemy.Figures.Set(eatenID, eatenFigure)
		}
	}
	check := f.team.CheckingCheck()
	undoMove()
	undoEating()
	if check {
		return true
	}
	return false
}

// positionOnBoard returns true if the coordinates are within the board, otherwise return false
func (f *Figure) positionOnBoard(pos *Position) bool {
	x, y := pos.Get()
	if x >= 1 && x <= 8 && y >= 1 && y <= 8 {
		return true
	}
	return false
}

// SetTeam set links to your team for current figure
func (f *Figure) SetTeam(team *Team) {
	f.team = team
}

// Move to new position and eat enemy figure if you need that
func (f *Figure) Move(pos *Position) {
	// pawn double move
	if f.GetName() == "pawn" {
		f.team.enemy.pawnDoubleMove.pawnTakeOnThePass(pos)
		f.team.pawnDoubleMove.pawnMakesDoubleMove(f, f.GetPosition(), pos)
	} else {
		f.team.pawnDoubleMove.clearPawnDoubleMove()
	}

	f.SetPosition(pos)
	f.setAlreadyMove(true)
	if f.team.enemy != nil && f.team.enemy.Figures.ExistsByPosition(pos) {
		// eat enemy figure
		err := f.team.enemy.Eating(pos)
		if err != nil {
			log.Println(err)
		}
	}

	f.transformPawnTOQueen(pos)

	// Debug
	//f.team.ShowBrokenFields()
	//f.team.enemy.ShowBrokenFields()
	//f.team.ShowPossibleMoves()
	//f.team.enemy.ShowPossibleMoves()
}

// transformPawnToQueen promote a pawn to a queen
func (f *Figure) transformPawnTOQueen(pos *Position) {
	if f.GetName() == "pawn" && (f.Y == 1 || f.Y == 8) {
		figureID := f.team.Figures.GetIndexByPosition(pos)
		// replace pawn to queen
		f.team.Figures.Set(figureID, NewQueen(pos, f.team))
		f.team.Figures.Get(figureID).setAlreadyMove(true)
	}
}

// SetPosition set Position to coords from argument
func (f *Figure) SetPosition(p *Position) {
	f.Position = p
}

// GetPosition return current Position of Figure
func (f *Figure) GetPosition() *Position {
	return f.Position
}
