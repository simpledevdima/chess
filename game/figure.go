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
	SetTeams(*Team, *Team)
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
	enemy       *Team
}

// GetName get name from figure
func (figureData *figureData) GetName() string {
	return figureData.Name
}

// SetName set name to the figure
func (figureData *figureData) SetName(name string) {
	figureData.Name = name
}

// IsAlreadyMove returns true if the figure has moved in the current match otherwise returns false
func (figureData *figureData) IsAlreadyMove() bool {
	return figureData.alreadyMove
}

// setAlreadyMove setting a value about moving a figure in the current match
func (figureData *figureData) setAlreadyMove(flag bool) {
	figureData.alreadyMove = flag
}

// kingOnTheBeatenFieldAfterMove returns true if your king is on the beaten square after the move otherwise return false
func (figureData *figureData) kingOnTheBeatenFieldAfterMove(x int, y int) bool {
	curX, curY := figureData.Position.Get()
	figureData.SetPosition(x, y)
	undoMove := func() {
		figureData.SetPosition(curX, curY)
	}
	undoEating := func() {}
	if figureData.enemy.FigureExist(x, y) {
		eatenID, _ := figureData.enemy.GetFigureID(x, y)
		eatenFigure := figureData.enemy.Figures[eatenID]
		delete(figureData.enemy.Figures, eatenID)
		undoEating = func() {
			figureData.enemy.Figures[eatenID] = eatenFigure
		}
	}
	check := figureData.team.CheckingCheck()
	undoMove()
	undoEating()
	if check {
		return true
	}
	return false
}

// coordsOnBoard returns true if the coordinates are within the board, otherwise return false
func (figureData *figureData) coordsOnBoard(x int, y int) bool {
	if x >= 1 && x <= 8 && y >= 1 && y <= 8 {
		return true
	}
	return false
}

// SetTeams set links to your team and enemy team for current figure
func (figureData *figureData) SetTeams(team *Team, enemy *Team) {
	figureData.team = team
	figureData.enemy = enemy
}

// MoveFigure to new coords and eat enemy figure if need that
func (figureData *figureData) MoveFigure(x int, y int) {
	figureData.SetPosition(x, y)
	figureData.setAlreadyMove(true)
	if figureData.enemy != nil && figureData.enemy.FigureExist(x, y) {
		// eat enemy figure
		err := figureData.enemy.Eating(x, y)
		if err != nil {
			log.Println(err)
		}
	}
}

// SetPosition set Position to coords from argument
func (figureData *figureData) SetPosition(x, y int) {
	figureData.Position.Set(x, y)
}

// GetPosition return current Position of Figure
func (figureData *figureData) GetPosition() (int, int) {
	return figureData.Position.Get()
}
