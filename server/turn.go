package server

import (
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
)

// newTurn returns a reference to the new turn structure
func newTurn() *turn {
	turn := &turn{}
	return turn
}

// turn data structure denoting turn queue
type turn struct {
	teamName game.TeamName
	server   *server
}

// now return teamName whose turn is currently active
func (t *turn) now() game.TeamName {
	return t.teamName
}

// getNowString return string with teamName whose turn is currently active
func (t *turn) getNowString() string {
	switch t.teamName {
	case game.White:
		return "white"
	case game.Black:
		return "black"
	}
	return "unknown team"
}

// setServer set link to the server
func (t *turn) setServer(server *server) {
	t.server = server
}

// setDefault set default values for a new game
func (t *turn) setDefault() {
	t.teamName = game.White

	// Debug
	t.server.board.White.ShowBrokenFields(t.server.board.White.GetBrokenFields())
	t.server.board.White.ShowPossibleMoves(t.server.board.White.GetPossibleMoves())
}

// change transfers the turn to the opposing team
func (t *turn) change() {
	if t.server.status.isPlay() {
		switch t.now() {
		case game.White:
			t.server.timers.white.stop()
			t.server.timers.white.reset()
			if t.server.board.Black.HavePossibleMove() {

				// Debug
				t.server.board.Black.ShowBrokenFields(t.server.board.Black.GetBrokenFields())
				t.server.board.Black.ShowPossibleMoves(t.server.board.Black.GetPossibleMoves())

				t.teamName = game.Black
				go t.server.timers.black.play()
				t.send(t.exportJSON())
			} else {
				if t.server.board.Black.CheckingCheck() {
					t.server.status.setOverCauseToWhite()
				} else {
					t.server.status.setOverCauseToStalemate()
				}
			}
		case game.Black:
			t.server.timers.black.stop()
			t.server.timers.black.reset()
			if t.server.board.White.HavePossibleMove() {

				// Debug
				t.server.board.White.ShowBrokenFields(t.server.board.White.GetBrokenFields())
				t.server.board.White.ShowPossibleMoves(t.server.board.White.GetPossibleMoves())

				t.teamName = game.White
				go t.server.timers.white.play()
				t.send(t.exportJSON())
			} else {
				if t.server.board.White.CheckingCheck() {
					t.server.status.setOverCauseToBlack()
				} else {
					t.server.status.setOverCauseToStalemate()
				}
			}
		}
	}
}

// exportJSON return data with current turn in JSON format
func (t *turn) exportJSON() []byte {
	request := &nrp.Simple{Post: "turn", Body: struct {
		White bool `json:"white,omitempty"`
		Black bool `json:"black,omitempty"`
	}{
		White: t.now() == game.White,
		Black: t.now() == game.Black,
	}}
	return request.Export()
}

// send data to broadcast
func (t *turn) send(dataJSON []byte) {
	t.server.broadcast <- dataJSON
}
