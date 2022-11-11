package server

import (
	"encoding/json"
	"github.com/skvdmt/chess/game"
	"log"
)

// turn data structure denoting turn queue
type turn struct {
	teamName game.TeamName
	server   *server
}

// now return teamName whose turn is currently active
func (turn *turn) now() game.TeamName {
	return turn.teamName
}

// getNowString return string with teamName whose turn is currently active
func (turn *turn) getNowString() string {
	switch turn.teamName {
	case game.White:
		return "white"
	case game.Black:
		return "black"
	}
	return "unknown team"
}

// setServer set link to the server
func (turn *turn) setServer(server *server) {
	turn.server = server
}

// setDefault set default values for a new game
func (turn *turn) setDefault() {
	turn.teamName = game.White
}

// change transfers the turn to the opposing team
func (turn *turn) change() {
	if turn.server.status.isPlay() {
		switch turn.now() {
		case game.White:
			turn.server.timers.white.stop()
			turn.server.timers.white.reset()
			if turn.server.board.Black.HavePossibleMove() {
				turn.teamName = game.Black
				go turn.server.timers.black.play()
				turn.send(turn.exportJSON())
			} else {
				if turn.server.board.Black.CheckingCheck() {
					turn.server.status.setOverCauseToWhite()
				} else {
					turn.server.status.setOverCauseToStalemate()
				}
			}
		case game.Black:
			turn.server.timers.black.stop()
			turn.server.timers.black.reset()
			if turn.server.board.White.HavePossibleMove() {
				turn.teamName = game.White
				go turn.server.timers.white.play()
				turn.send(turn.exportJSON())
			} else {
				if turn.server.board.White.CheckingCheck() {
					turn.server.status.setOverCauseToBlack()
				} else {
					turn.server.status.setOverCauseToStalemate()
				}
			}
		}
	}
}

// exportJSON return data with current turn in JSON format
func (turn *turn) exportJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		Turn string `json:"turn"`
	}{
		Turn: turn.getNowString(),
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// send data to broadcast
func (turn *turn) send(dataJSON []byte) {
	turn.server.broadcast <- dataJSON
}
