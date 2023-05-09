package game

import (
	"github.com/simpledevdima/nrp"
	"log"
)

// NewBoard returns a link to a new board with the three teams created
func NewBoard() *Board {
	return &Board{
		White:      NewTeam(White),
		Black:      NewTeam(Black),
		Spectators: NewTeam(Spectators),
	}
}

// Board data type with information about the board and its control
type Board struct {
	White      *Team `json:"white"`
	Black      *Team `json:"black"`
	Spectators *Team `json:"-"`
}

// NewBoard making a new board
func (b *Board) NewBoard() {
	err := b.White.setStartPosition()
	if err != nil {
		log.Println(err)
	}
	err = b.Black.setStartPosition()
	if err != nil {
		log.Println(err)
	}
}

// ExportJSON getting data about all the Figures on the board in JSON format
func (b *Board) ExportJSON() []byte {
	request := nrp.Simple{Post: "board", Body: b}
	return request.Export()
}
