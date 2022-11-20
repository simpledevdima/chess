package game

import (
	"log"
)

// figureData information about each figure
type figureData struct {
	Name        string `json:"name"`
	*Position   `json:"position"`
	alreadyMove bool
	team        *Team
}

// GetName get name from figure
func (f *figureData) GetName() string {
	return f.Name
}

// SetName set name to the figure
func (f *figureData) SetName(name string) {
	f.Name = name
}

// IsAlreadyMove returns true if the figure has moved in the current match otherwise returns false
func (f *figureData) IsAlreadyMove() bool {
	return f.alreadyMove
}

// setAlreadyMove setting a value about moving a figure in the current match
func (f *figureData) setAlreadyMove(flag bool) {
	f.alreadyMove = flag
}

// kingOnTheBeatenFieldAfterMove returns true if your king is on the beaten square after the move otherwise return false
func (f *figureData) kingOnTheBeatenFieldAfterMove(pos *Position) bool {
	curPos := f.Position
	f.SetPosition(pos)
	undoMove := func() {
		f.SetPosition(curPos)
	}
	undoEating := func() {}
	if f.team.enemy.Figures.ExistsByPosition(pos) {
		eatenID := f.team.enemy.Figures.GetIndexByPosition(pos)
		eatenFigure := f.team.enemy.Figures.Get(eatenID)
		f.team.enemy.Figures.RemoveByIndex(eatenID)
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
func (f *figureData) positionOnBoard(pos *Position) bool {
	x, y := pos.Get()
	if x >= 1 && x <= 8 && y >= 1 && y <= 8 {
		return true
	}
	return false
}

// SetTeam set links to your team for current figure
func (f *figureData) SetTeam(team *Team) {
	f.team = team
}

// MoveFigure to new coords and eat enemy figure if you need that
func (f *figureData) MoveFigure(position *Position) {
	f.SetPosition(position)
	f.setAlreadyMove(true)
	if f.team.enemy != nil && f.team.enemy.Figures.ExistsByPosition(position) {
		// eat enemy figure
		err := f.team.enemy.Eating(position)
		if err != nil {
			log.Println(err)
		}
	}

	// Debug
	//f.team.ShowBrokenFields()
	//f.team.enemy.ShowBrokenFields()
	//f.team.ShowPossibleMoves()
	//f.team.enemy.ShowPossibleMoves()
}

// SetPosition set Position to coords from argument
func (f *figureData) SetPosition(p *Position) {
	f.Position = p
}

// GetPosition return current Position of Figure
func (f *figureData) GetPosition() *Position {
	return f.Position
}
