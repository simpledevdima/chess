package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func NewTeam(n TeamName) *Team {
	return &Team{Name: n}
}

// Team data type for white or black team
type Team struct {
	Name           TeamName
	Figures        Figures `json:"figures"`
	Eaten          Figures `json:"eaten"`
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
			if ok, _ := figure.Validation(position); ok {
				return true
			}
		}
	}
	return false
}

// CheckingCheck returns true if the king is on a beaten field otherwise returns false
func (t *Team) CheckingCheck() bool {
	king := t.Figures.GetByName("king")
	kingPos := king.GetPosition()
	for _, figure := range t.enemy.Figures {
		for _, figPos := range figure.detectionOfBrokenFields() {
			if *figPos == *kingPos {
				return true
			}
		}
	}
	return false
}

// Eating figure on x, y coords move its figure from Figures map to Eaten map
func (t *Team) Eating(eatPos *Position) error {
	for id, figure := range t.Figures {
		figPos := figure.GetPosition()
		if *figPos == *eatPos {
			t.Eaten[id] = figure
			delete(t.Figures, id)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("we cant eat figure because no figure in: %vx%v coords", eatPos.X, eatPos.Y))
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
	t.MakeFigures()
	// paws
	for x := 1; x <= 8; x++ {
		t.Figures[FigureIndex(x)] = NewPawn(NewPosition(x, pawnLine), t)
	}
	// rooks
	t.Figures[9] = NewRook(NewPosition(1, figuresLine), t)
	t.Figures[16] = NewRook(NewPosition(8, figuresLine), t)
	// knights
	t.Figures[10] = NewKnight(NewPosition(2, figuresLine), t)
	t.Figures[15] = NewKnight(NewPosition(7, figuresLine), t)
	// bishops
	t.Figures[11] = NewBishop(NewPosition(3, figuresLine), t)
	t.Figures[14] = NewBishop(NewPosition(6, figuresLine), t)
	// king
	t.Figures[12] = NewKing(NewPosition(5, figuresLine), t)
	// queen
	t.Figures[13] = NewQueen(NewPosition(4, figuresLine), t)
	return nil
}

// MakeFigures remake Figures and Eaten map
func (t *Team) MakeFigures() {
	t.Figures = make(Figures)
	t.Eaten = make(Figures)
}

// ImportFigures sets the data received in JSON format from the argument to the command shapes
func (t *Team) ImportFigures(figuresJSON []byte) {
	t.MakeFigures()
	var figures map[FigureIndex]struct {
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
		pos := NewPosition(figure.Position.X, figure.Position.Y)
		switch figure.Name {
		case "pawn":
			t.Figures[index] = NewPawn(pos, t)
		case "knight":
			t.Figures[index] = NewKnight(pos, t)
		case "bishop":
			t.Figures[index] = NewBishop(pos, t)
		case "rook":
			t.Figures[index] = NewRook(pos, t)
		case "queen":
			t.Figures[index] = NewQueen(pos, t)
		case "king":
			t.Figures[index] = NewKing(pos, t)
		}
	}
}

// PossibleMoves data type with possible moves of pieces
type PossibleMoves map[FigureIndex][]*Position

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
			x, y := figure.GetPosition().Get()
			fmt.Println(index, figure.GetName(), x, y, fields)
		}
	}
}

// ShowBrokenFields displays the squares that beat the figures of the team
func (t *Team) ShowBrokenFields() {
	fmt.Println("Team:", t.Name.String())
	for index, figure := range t.Figures {
		fields := figure.detectionOfBrokenFields()
		pos := figure.GetPosition()
		fmt.Println(index, figure.GetName(), pos.X, pos.Y, fields)
	}
}
