package game

import (
	"encoding/json"
	"log"
)

// Board data type with information about the board and its control
type Board struct {
	boardData      `json:"board"`
	Spectators     *Team `json:"-"`
	pawnDoubleMove pawnDoubleMove
}

// boardData
type boardData struct {
	White *Team `json:"white"`
	Black *Team `json:"black"`
}

// setTeam
func (board *Board) setTeam(teamName TeamName, team *Team) {
	switch teamName {
	case White:
		board.White = team
	case Black:
		board.Black = team
	case Spectators:
		board.Spectators = team
	}
}

// MakeTeams setup teams on the Board
func (board *Board) MakeTeams() {
	board.setTeam(White, &Team{Name: White})
	board.setTeam(Black, &Team{Name: Black})
	board.setTeam(Spectators, &Team{Name: Spectators})
}

// NewBoard making a new board
func (board *Board) NewBoard() {
	err := board.White.setStartPosition()
	if err != nil {
		log.Println(err)
	}
	err = board.Black.setStartPosition()
	if err != nil {
		log.Println(err)
	}
	// set link to teams and board for all Figures
	for _, figure := range board.White.Figures {
		figure.SetTeams(board.White, board.Black)
		figure.setBoard(board)
	}
	for _, figure := range board.Black.Figures {
		figure.SetTeams(board.Black, board.White)
		figure.setBoard(board)
	}
}

// ExportJSON getting data about all the Figures on the board in JSON format
func (board *Board) ExportJSON() []byte {
	dataJSON, err := json.Marshal(board)
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}
