package game

import (
	"log"
)

// Figure a set of methods for any chess figure (king, queen, rook, knight, bishop, pawn)
type Figure interface {
	GetName() string
	SetName(string)
	SetPosition(int, int)
	GetPosition() (int, int)
	Move(int, int)
	MoveFigure(int, int)
	Validation(int, int) (bool, string)
	SetTeam(*Team)
	coordsOnBoard(int, int) bool
	kingOnTheBeatenFieldAfterMove(int, int) bool
	detectionOfBrokenFields() []Position
	DetectionOfPossibleMove() []Position
	IsAlreadyMove() bool
	setAlreadyMove(bool)
}

// figureData information about each figure
type figureData struct {
	Name        string `json:"name"`
	Position    `json:"position"`
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
func (f *figureData) kingOnTheBeatenFieldAfterMove(x int, y int) bool {
	curX, curY := f.Position.Get()
	f.SetPosition(x, y)
	undoMove := func() {
		f.SetPosition(curX, curY)
	}
	undoEating := func() {}
	if f.team.enemy.FigureExist(x, y) {
		eatenID, _ := f.team.enemy.GetFigureID(x, y)
		eatenFigure := f.team.enemy.Figures[eatenID]
		delete(f.team.enemy.Figures, eatenID)
		undoEating = func() {
			f.team.enemy.Figures[eatenID] = eatenFigure
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

// coordsOnBoard returns true if the coordinates are within the board, otherwise return false
func (f *figureData) coordsOnBoard(x int, y int) bool {
	if x >= 1 && x <= 8 && y >= 1 && y <= 8 {
		return true
	}
	return false
}

// SetTeam set links to your team for current figure
func (f *figureData) SetTeam(team *Team) {
	f.team = team
}

// MoveFigure to new coords and eat enemy figure if need that
func (f *figureData) MoveFigure(x int, y int) {
	f.SetPosition(x, y)
	f.setAlreadyMove(true)
	if f.team.enemy != nil && f.team.enemy.FigureExist(x, y) {
		// eat enemy figure
		err := f.team.enemy.Eating(x, y)
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
func (f *figureData) SetPosition(x, y int) {
	f.Position.Set(x, y)
}

// GetPosition return current Position of Figure
func (f *figureData) GetPosition() (int, int) {
	return f.Position.Get()
}
