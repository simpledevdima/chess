package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Team data type for white or black team
type Team struct {
	Name           TeamName
	Figures        map[int]Figure `json:"figures"`
	Eaten          map[int]Figure
	enemy          *Team
	pawnDoubleMove pawnDoubleMove
}

// SetName exported method of setting the command name by the string from the argument
func (team *Team) SetName(teamName string) {
	switch teamName {
	case "white":
		team.Name = White
	case "black":
		team.Name = Black
	case "spectators":
		team.Name = Spectators
	default:
		log.Println(errors.New(fmt.Sprintf("unknown team name: %s", teamName)))
	}
}

// SetEnemy set link to enemy of team
func (team *Team) SetEnemy(enemy *Team) {
	team.enemy = enemy
}

// HavePossibleMove return true if team can make move
func (team *Team) HavePossibleMove() bool {
	for _, figure := range team.Figures {
		for _, position := range figure.DetectionOfPossibleMove() {
			if ok, _ := figure.Validation(position.X, position.Y); ok {
				return true
			}
		}
	}
	return false
}

// CheckingCheck returns true if the king is on a beaten field otherwise returns false
func (team *Team) CheckingCheck() bool {
	kingID, err := team.getFigureIDByName("king")
	if err != nil {
		log.Println(err)
	}
	x, y := team.Figures[kingID].GetPosition()
	for _, figure := range team.enemy.Figures {
		for _, position := range figure.detectionOfBrokenFields() {
			if position.X == x && position.Y == y {
				return true
			}
		}
	}
	return false
}

// getFigureIDByName get figure by name and return ID and error
func (team *Team) getFigureIDByName(name string) (int, error) {
	for id, figure := range team.Figures {
		if strings.ToLower(reflect.TypeOf(figure).Elem().Name()) == name && figure.GetName() == name {
			return id, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("figure \"%s\" not forund", name))
}

// GetFigureID return ID and error by coords
func (team *Team) GetFigureID(x int, y int) (int, error) {
	for id, figure := range team.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			return id, nil
		}
	}
	return 0, errors.New("figure not exist")
}

// Eating figure on x, y coords move its figure from Figures map to Eaten map
func (team *Team) Eating(x int, y int) error {
	for id, figure := range team.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			team.Eaten[id] = figure
			delete(team.Figures, id)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("we cant eat figure because no figure in: %vx%v coords", x, y))
}

// FigureExist check Figures of the team and return true if figure exist on take arguments coords else return false
func (team *Team) FigureExist(x int, y int) bool {
	for _, figure := range team.Figures {
		figX, figY := figure.GetPosition()
		if figX == x && figY == y {
			return true
		}
	}
	return false
}

// setStartPosition method setup start team positions for all Figures
func (team *Team) setStartPosition() error {
	var FiguresLine int
	var pawnLine int
	switch team.Name {
	case White:
		FiguresLine = 1
		pawnLine = 2
	case Black:
		FiguresLine = 8
		pawnLine = 7
	default:
		return errors.New("undefined team name")
	}
	team.ClearFigures()
	// paws
	for i := 1; i <= 8; i++ {
		team.Figures[i] = &Pawn{}
		team.Figures[i].SetName("pawn")
		team.Figures[i].SetPosition(i, pawnLine)
		team.Figures[i] = &Pawn{}
		team.Figures[i].SetName("pawn")
		team.Figures[i].SetPosition(i, pawnLine)
	}
	// rooks
	team.Figures[9] = &Rook{}
	team.Figures[9].SetName("rook")
	team.Figures[9].SetPosition(1, FiguresLine)
	team.Figures[16] = &Rook{}
	team.Figures[16].SetName("rook")
	team.Figures[16].SetPosition(8, FiguresLine)
	// knights
	team.Figures[10] = &Knight{}
	team.Figures[10].SetName("knight")
	team.Figures[10].SetPosition(2, FiguresLine)
	team.Figures[15] = &Knight{}
	team.Figures[15].SetName("knight")
	team.Figures[15].SetPosition(7, FiguresLine)
	// bishops
	team.Figures[11] = &Bishop{}
	team.Figures[11].SetName("bishop")
	team.Figures[11].SetPosition(3, FiguresLine)
	team.Figures[14] = &Bishop{}
	team.Figures[14].SetName("bishop")
	team.Figures[14].SetPosition(6, FiguresLine)
	// king
	team.Figures[12] = &King{}
	team.Figures[12].SetName("king")
	team.Figures[12].SetPosition(5, FiguresLine)
	// queen
	team.Figures[13] = &Queen{}
	team.Figures[13].SetName("queen")
	team.Figures[13].SetPosition(4, FiguresLine)

	// set links to team and enemy for all Figures
	for _, figure := range team.Figures {
		figure.SetTeams(team, team.enemy)
	}

	team.GetPossibleMoves()

	return nil
}

// ClearFigures remake Figures and Eaten map
func (team *Team) ClearFigures() {
	team.Figures = make(map[int]Figure)
	team.Eaten = make(map[int]Figure)
}

// ImportFigures sets the data received in JSON format from the argument to the command shapes
func (team *Team) ImportFigures(figuresJSON []byte) {
	team.ClearFigures()
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
			team.Figures[index] = &Pawn{}
		case "knight":
			team.Figures[index] = &Knight{}
		case "bishop":
			team.Figures[index] = &Bishop{}
		case "rook":
			team.Figures[index] = &Rook{}
		case "queen":
			team.Figures[index] = &Queen{}
		case "king":
			team.Figures[index] = &King{}
		}
		team.Figures[index].SetPosition(figure.Position.X, figure.Position.Y)
		team.Figures[index].SetTeams(team, team.enemy)
		team.Figures[index].SetName(figure.Name)
	}
}

// PossibleMoves data type with possible moves of pieces
type PossibleMoves map[int][]Position

// GetPossibleMoves returns a map with the keys of the team's shapes and the slices of coordinates that those shapes can make
func (team *Team) GetPossibleMoves() PossibleMoves {
	possibleMoves := make(PossibleMoves)
	for index, figure := range team.Figures {
		moves := figure.DetectionOfPossibleMove()
		if len(moves) > 0 {
			possibleMoves[index] = moves
		}
	}
	return possibleMoves
}
