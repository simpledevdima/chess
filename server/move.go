package server

import (
	"encoding/json"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
	"log"
)

func newMove(c *client) *move {
	m := &move{}
	m.setClient(c)
	return m
}

// move data type containing the processing of the movement of the figure
type move struct {
	From struct {
		*game.Position `json:"position"`
	} `json:"from"`
	To struct {
		*game.Position `json:"position"`
	} `json:"to"`
	client *client
}

// setClient set a link to the client that makes the move
func (m *move) setClient(client *client) {
	m.client = client
}

// isValid returns true and an empty string if it is possible to make a move otherwise returns false and a string with a value on which a move is not possible
func (m *move) isValid() (bool, string) {
	if m.client.server.status.isPlay() {
		if m.client.server.turn.now() == m.client.team.Name {
			if m.client.team.Figures.ExistsByPosition(m.From.Position) {
				if ok, cause := m.client.team.Figures.GetByPosition(m.From.Position).Validation(m.To.Position); ok {
					return true, ""
				} else {
					return false, cause
				}
			} else {
				return false, "wrong figure"
			}
		} else {
			return false, "now not your move"
		}
	} else {
		return false, "game are stopped"
	}
}

// exec executes the current move and sends the data to the broadcast
func (m *move) exec() {
	figure := m.client.team.Figures.GetByPosition(m.From.Position)
	if m.isCastling(figure) {
		m.makeRookMoveInCastling()
	}
	figure.Move(m.To.Position)
	event := nrp.Simple{Post: "move", Body: &m}
	m.client.server.broadcast <- event.Export()
}

// isCastling returns true if the move is castling otherwise returns false
func (m *move) isCastling(figure game.Figure) bool {
	if figure.GetName() == "king" && !figure.IsAlreadyMove() && (m.To.Position.X == 3 || m.To.Position.X == 7) {
		return true
	}
	return false
}

// makeCastling creates a rook move and makes it
func (m *move) makeRookMoveInCastling() {
	moveRook := newMove(m.client)
	moveRook.From.Position.Y = m.From.Position.Y
	moveRook.To.Position.Y = m.To.Position.Y
	switch m.To.Position.X {
	case 3:
		moveRook.From.Position.X = 1
		moveRook.To.Position.X = 4
	case 7:
		moveRook.From.Position.X = 8
		moveRook.To.Position.X = 6
	}
	moveRook.exec()
}

// exportJSON get data of current type in JSON format
func (m *move) exportJSON() []byte {
	dataJSON, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}
