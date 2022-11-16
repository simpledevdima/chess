package server

import (
	"encoding/json"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
	"log"
)

func NewMove(c *client) *move {
	m := &move{}
	m.setClient(c)
	return m
}

// move data type containing the processing of the movement of the figure
type move struct {
	From struct {
		game.Position `json:"position"`
	} `json:"from"`
	To struct {
		game.Position `json:"position"`
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
			if figureID, err := m.client.team.GetFigureID(m.From.Position.X, m.From.Position.Y); err == nil {
				if ok, cause := m.client.team.Figures[figureID].Validation(m.To.Position.X, m.To.Position.Y); ok {
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
	figureID, err := m.client.team.GetFigureID(m.From.Position.X, m.From.Position.Y)
	if err != nil {
		log.Println(err)
	}
	// castling
	if !m.client.team.Figures[figureID].IsAlreadyMove() && m.From.Position.X == 5 && (m.From.Position.Y == 1 || m.From.Position.Y == 8) && (m.To.Position.X == 3 || m.To.Position.X == 7) {
		var m2 move
		m2.From.Position.Y = m.From.Position.Y
		m2.To.Position.Y = m.To.Position.Y
		switch m.To.Position.X {
		case 3:
			m2.From.Position.X = 1
			m2.To.Position.X = 4
		case 7:
			m2.From.Position.X = 8
			m2.To.Position.X = 6
		}
		m2.setClient(m.client)
		m2.exec()
	}

	m.client.team.Figures[figureID].Move(m.To.Position.X, m.To.Position.Y)
	event := nrp.Simple{Post: "move", Body: &m}
	m.client.server.broadcast <- event.Export()
}

// exportJSON get data of current type in JSON format
func (m *move) exportJSON() []byte {
	dataJSON, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}
