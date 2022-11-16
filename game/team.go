package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func NewTeam(n TeamName) *Team {
	return &Team{Name: n}
}

// Team data type for white or black team
type Team struct {
	Name           TeamName
	Figures        map[int]Figure `json:"figures"`
	Eaten          map[int]Figure
	enemy          *Team
	pawnDoubleMove pawnDoubleMove // taking on the pass
}

// SetName exported method of setting the command name by the string from the argument
func (t *Team) SetName(teamName string) {
	switch teamName {
	case "white":
		t.Name = White
	case "black":
		t.Name = Black
	case "spectators":
		t.Name = Spectators
	default:
		log.Println(errors.New(fmt.Sprintf("unknown team name: %s", teamName)))
	}
}

// SetEnemy set link to enemy of team
func (t *Team) SetEnemy(enemy *Team) {
	t.enemy = enemy
}

// HavePossibleMove return true if team can make move
func (t *Team) HavePossibleMove() bool {
	for _, figure := range t.Figures {
		for _, position := range figure.DetectionOfPossibleMove() {
			if ok, _ := figure.Validation(position.X, position.Y); ok {
				return true
			}
		}
	}
	return false
}

// CheckingCheck returns true if the king is on a beaten field otherwise returns false
func (t *Team) CheckingCheck() bool {
	kingID, err := t.getFigureIDByName("king")
	if err != nil {
		log.Println(err)
	}
	x, y := t.Figures[kingID].GetPosition()
	for _, figure := range t.enemy.Figures {
		for _, position := range figure.detectionOfBrokenFields() {
			if position.X == x && position.Y == y {
				return true
			}
		}
	}
	return false
}

// getFigureIDByName get figure by name and return ID and error
func (t *Team) getFigureIDByName(name string) (int, error) {
	for id, figure := range t.Figures {
		if strings.ToLower(reflect.TypeOf(figure).Elem().Name()) == name && figure.GetName() == name {
			return id, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("figure \"%s\" not forund", name))
}

// GetFigureID return ID and error by coords
func (t *Team) GetFigureID(x int, y int) (int, error) {
	for id, figure := range t.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			return id, nil
		}
	}
	return 0, errors.New("figure not exist")
}

// GetFigureByCoords return link to the figure by coords if found in team
func (t *Team) GetFigureByCoords(x, y int) Figure {
	for _, figure := range t.Figures {
		fX, fY := figure.GetPosition()
		if fX == x && fY == y {
			return figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure by cords %d x %d not found in %s team\n", x, y, t.Name.String())))
	return nil
}

// Eating figure on x, y coords move its figure from Figures map to Eaten map
func (t *Team) Eating(x int, y int) error {
	for id, figure := range t.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			t.Eaten[id] = figure
			delete(t.Figures, id)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("we cant eat figure because no figure in: %vx%v coords", x, y))
}

// FigureExist check Figures of the team and return true if figure exist on take arguments coords else return false
func (t *Team) FigureExist(x int, y int) bool {
	for _, figure := range t.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			return true
		}
	}
	return false
}

// setStartPosition method setup start team positions for all Figures
func (t *Team) setStartPosition() error {
	var figuresLine int
	var pawnLine int
	switch t.Name {
	case White:
		figuresLine = 1
		pawnLine = 2
	case Black:
		figuresLine = 8
		pawnLine = 7
	default:
		return errors.New("undefined team name")
	}
	t.ClearFigures()
	// paws
	for x := 1; x <= 8; x++ {
		t.Figures[x] = NewPawn(x, pawnLine, t)
	}
	// rooks
	t.Figures[9] = NewRook(1, figuresLine, t)
	t.Figures[16] = NewRook(8, figuresLine, t)
	// knights
	t.Figures[10] = NewKnight(2, figuresLine, t)
	t.Figures[15] = NewKnight(7, figuresLine, t)
	// bishops
	t.Figures[11] = NewBishop(3, figuresLine, t)
	t.Figures[14] = NewBishop(6, figuresLine, t)
	// king
	t.Figures[12] = NewKing(5, figuresLine, t)
	// queen
	t.Figures[13] = NewQueen(4, figuresLine, t)
	return nil
}

// ClearFigures remake Figures and Eaten map
func (t *Team) ClearFigures() {
	t.Figures = make(map[int]Figure)
	t.Eaten = make(map[int]Figure)
}

// ImportFigures sets the data received in JSON format from the argument to the command shapes
func (t *Team) ImportFigures(figuresJSON []byte) {
	t.ClearFigures()
	var figures map[int]struct {
		Name     string `json:"name"`
		Position struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"position"`
	}
	err := json.Unmarshal(figuresJSON, &figures)
	if err != nil {
		log.Println(err)
	}
	for index, figure := range figures {
		switch figure.Name {
		case "pawn":
			t.Figures[index] = NewPawn(figure.Position.X, figure.Position.Y, t)
		case "knight":
			t.Figures[index] = NewKnight(figure.Position.X, figure.Position.Y, t)
		case "bishop":
			t.Figures[index] = NewBishop(figure.Position.X, figure.Position.Y, t)
		case "rook":
			t.Figures[index] = NewRook(figure.Position.X, figure.Position.Y, t)
		case "queen":
			t.Figures[index] = NewQueen(figure.Position.X, figure.Position.Y, t)
		case "king":
			t.Figures[index] = NewKing(figure.Position.X, figure.Position.Y, t)
		}
	}
}

// PossibleMoves data type with possible moves of pieces
type PossibleMoves map[int][]Position

// GetPossibleMoves returns a map with the keys of the team's shapes and the slices of coordinates that those shapes can make
func (t *Team) GetPossibleMoves() PossibleMoves {
	possibleMoves := make(PossibleMoves)
	for index, figure := range t.Figures {
		moves := figure.DetectionOfPossibleMove()
		if len(moves) > 0 {
			possibleMoves[index] = moves
		}
	}
	return possibleMoves
}

// ShowPossibleMoves displays the possible moves of each piece of the team
func (t *Team) ShowPossibleMoves() {
	fmt.Println("Team:", t.Name.String())
	for index, figure := range t.Figures {
		fields := figure.DetectionOfPossibleMove()
		if len(fields) > 0 {
			x, y := figure.GetPosition()
			fmt.Println(index, figure.GetName(), x, y, fields)
		}
	}
}

// ShowBrokenFields displays the squares that beat the figures of the team
func (t *Team) ShowBrokenFields() {
	fmt.Println("Team:", t.Name.String())
	for index, figure := range t.Figures {
		fields := figure.detectionOfBrokenFields()
		x, y := figure.GetPosition()
		fmt.Println(index, figure.GetName(), x, y, fields)
	}
}
