package game

import (
	"github.com/skvdmt/nrp"
	"log"
)

// Board data type with information about the board and its control
type Board struct {
	White      *Team `json:"white"`
	Black      *Team `json:"black"`
	Spectators *Team `json:"-"`
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
}

// ExportJSON getting data about all the Figures on the board in JSON format
func (board *Board) ExportJSON() []byte {
	request := nrp.Simple{Post: "board", Body: board}
	return request.Export()
}
