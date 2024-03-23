package server

import (
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
)

// newStatus returns a reference to the new status structure
func newStatus() *status {
	return &status{}
}

// status data structure storing various states of the game process
type status struct {
	play      bool
	waitCause waitCause
	over      bool
	overCause overCause
	server    *server
}

// waitCause reason for waiting
type waitCause int

const (
	waitBothPlayers waitCause = iota
	waitWhitePlayer
	waitBlackPlayer
	allPayersReady
)

// overCause termination reason
type overCause int

const (
	winnerNotSet overCause = iota
	winnerIsWhite
	winnerIsBlack
	stalemate
	drawGame
)

// setServer sets a link to the server
func (s *status) setServer(server *server) {
	s.server = server
}

// isPlay returns the play state
func (s *status) isPlay() bool {
	return s.play
}

// isOver returns the completion status
func (s *status) isOver() bool {
	return s.over
}

// changePlayAndSend changes the play state and sends data about it to the broadcast
func (s *status) changePlayAndSend() {
	s.changePlay()
	s.send(s.exportPlayJSON())
}

// changePlay changes playback state
func (s *status) changePlay() {
	if s.waitCause == allPayersReady {
		s.play = true
		if !s.isOver() {
			s.server.play()
		}
	} else {
		if s.isPlay() {
			s.play = false
			s.server.stop()
		}
	}
}

// setWaitBothPlayers sets the reason for the wait to be both players' wait
func (s *status) setWaitBothPlayers() {
	s.waitCause = waitBothPlayers
	s.changePlayAndSend()
}

// setWaitWhitePlayer sets the cause of the wait to the value of waiting for a white team player
func (s *status) setWaitWhitePlayer() {
	s.waitCause = waitWhitePlayer
	s.changePlayAndSend()
}

// setWaitBlackPlayer sets the reason for waiting to the value of waiting for a black team player
func (s *status) setWaitBlackPlayer() {
	s.waitCause = waitBlackPlayer
	s.changePlayAndSend()
}

// setAllPlayersReady sets wait reason to all players ready
func (s *status) setAllPlayersReady() {
	s.waitCause = allPayersReady
	s.changePlayAndSend()
	s.server.setClientsEnemyLinks()
}

// isWaitBothPlayers returns true if both players are not present and the wait reason is not set to this value otherwise returns false
func (s *status) isWaitBothPlayers() bool {
	if !s.server.clientExistsByTeamName(game.White) &&
		!s.server.clientExistsByTeamName(game.Black) &&
		s.waitCause != waitBothPlayers {
		return true
	}
	return false
}

// isWaitWhitePlayer returns true if only the white team player is not present and the wait reason is not set to this value otherwise returns false
func (s *status) isWaitWhitePlayer() bool {
	if !s.server.clientExistsByTeamName(game.White) &&
		s.server.clientExistsByTeamName(game.Black) &&
		s.waitCause != waitWhitePlayer {
		return true
	}
	return false
}

// isWaitBlackPlayer returns true if only the black team player is not present and the wait reason is not set to this value otherwise returns false
func (s *status) isWaitBlackPlayer() bool {
	if s.server.clientExistsByTeamName(game.White) &&
		!s.server.clientExistsByTeamName(game.Black) &&
		s.waitCause != waitBlackPlayer {
		return true
	}
	return false
}

// isWaitBothPlayers returns true if both players are ready and the wait reason is not set to this value otherwise returns false
func (s *status) isAllPlayersReady() bool {
	if s.server.clientExistsByTeamName(game.White) &&
		s.server.clientExistsByTeamName(game.Black) &&
		s.waitCause != allPayersReady {
		return true
	}
	return false
}

// changeCausePlay change cause and server states depending on connecting and disconnecting clients
func (s *status) changeCausePlay() {
	if s.isWaitBothPlayers() {
		s.setWaitBothPlayers()
	} else if s.isWaitWhitePlayer() {
		s.setWaitWhitePlayer()
	} else if s.isWaitBlackPlayer() {
		s.setWaitBlackPlayer()
	} else if s.isAllPlayersReady() {
		s.setAllPlayersReady()
	}
}

// setOverCauseToWhite sets the reason for the end of the game to white win
func (s *status) setOverCauseToWhite() {
	s.overCause = winnerIsWhite
	s.changeOver()
}

// setOverCauseToBlack sets the reason for the end of the game to black win
func (s *status) setOverCauseToBlack() {
	s.overCause = winnerIsBlack
	s.changeOver()
}

// setOverCauseToStalemate sets the reason for ending the game to stalemate
func (s *status) setOverCauseToStalemate() {
	s.overCause = stalemate
	s.changeOver()
}

// setOverCauseToDraw sets game end reason to draw
func (s *status) setOverCauseToDraw() {
	s.overCause = drawGame
	s.changeOver()
}

// changeOver changes the value of the termination state to true if the values of the termination reason contribute to it and sends the data to the broadcast
func (s *status) changeOver() {
	if s.overCause != winnerNotSet {
		s.over = true
	}
	s.send(s.exportOverJSON())
}

// resetOver resets the termination value and termination reason to the default value
func (s *status) resetOver() {
	s.overCause = winnerNotSet
	s.over = false
}

// getCauseOver returns the string value of the termination reason
func (s *status) getCauseOver() string {
	switch s.overCause {
	case winnerIsWhite:
		return "white win"
	case winnerIsBlack:
		return "black win"
	case stalemate:
		return "stalemate"
	case drawGame:
		return "draw"
	}
	return ""
}

// exportOverJSON returns completion status and exit reason in JSON format
func (s *status) exportOverJSON() []byte {
	request := nrp.Simple{Post: "game_over", Body: struct {
		Over  bool   `json:"over"`
		Cause string `json:"cause,omitempty"`
	}{
		Over:  s.over,
		Cause: s.getCauseOver(),
	}}
	return request.Export()
}

// getCausePlay returns the string value of the replay reason
func (s *status) getCausePlay() string {
	switch s.waitCause {
	case waitBothPlayers:
		return "wait both players"
	case waitWhitePlayer:
		return "wait white player"
	case waitBlackPlayer:
		return "wait black player"
	}
	return ""
}

// exportPlayJSON returns the playback status and reason for playback in JSON format
func (s *status) exportPlayJSON() []byte {
	request := nrp.Simple{Post: "game_play", Body: struct {
		Play  bool   `json:"play"`
		Cause string `json:"cause,omitempty"`
	}{
		Play:  s.play,
		Cause: s.getCausePlay(),
	}}
	return request.Export()
}

// send data to broadcast
func (s *status) send(exportJSON []byte) {
	s.server.broadcast <- exportJSON
}
