package server

import (
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
)

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
func (status *status) setServer(server *server) {
	status.server = server
}

// isPlay returns the play state
func (status *status) isPlay() bool {
	return status.play
}

// isOver returns the completion status
func (status *status) isOver() bool {
	return status.over
}

// changePlayAndSend changes the play state and sends data about it to the broadcast
func (status *status) changePlayAndSend() {
	status.changePlay()
	status.send(status.exportPlayJSON())
}

// changePlay changes playback state
func (status *status) changePlay() {
	if status.waitCause == allPayersReady {
		status.play = true
		if !status.isOver() {
			status.server.play()
		}
	} else {
		if status.isPlay() {
			status.play = false
			status.server.stop()
		}
	}
}

// setWaitBothPlayers sets the reason for the wait to be both players' wait
func (status *status) setWaitBothPlayers() {
	status.waitCause = waitBothPlayers
	status.changePlayAndSend()
}

// setWaitWhitePlayer sets the cause of the wait to the value of waiting for a white team player
func (status *status) setWaitWhitePlayer() {
	status.waitCause = waitWhitePlayer
	status.changePlayAndSend()
}

// setWaitBlackPlayer sets the reason for waiting to the value of waiting for a black team player
func (status *status) setWaitBlackPlayer() {
	status.waitCause = waitBlackPlayer
	status.changePlayAndSend()
}

// setAllPlayersReady sets wait reason to all players ready
func (status *status) setAllPlayersReady() {
	status.waitCause = allPayersReady
	status.changePlayAndSend()
	status.server.setClientsEnemyLinks()
}

// isWaitBothPlayers returns true if both players are not present and the wait reason is not set to this value otherwise returns false
func (status *status) isWaitBothPlayers() bool {
	if !status.server.clientExistsByTeamName(game.White) &&
		!status.server.clientExistsByTeamName(game.Black) &&
		status.waitCause != waitBothPlayers {
		return true
	}
	return false
}

// isWaitWhitePlayer returns true if only the white team player is not present and the wait reason is not set to this value otherwise returns false
func (status *status) isWaitWhitePlayer() bool {
	if !status.server.clientExistsByTeamName(game.White) &&
		status.server.clientExistsByTeamName(game.Black) &&
		status.waitCause != waitWhitePlayer {
		return true
	}
	return false
}

// isWaitBlackPlayer returns true if only the black team player is not present and the wait reason is not set to this value otherwise returns false
func (status *status) isWaitBlackPlayer() bool {
	if status.server.clientExistsByTeamName(game.White) &&
		!status.server.clientExistsByTeamName(game.Black) &&
		status.waitCause != waitBlackPlayer {
		return true
	}
	return false
}

// isWaitBothPlayers returns true if both players are ready and the wait reason is not set to this value otherwise returns false
func (status *status) isAllPlayersReady() bool {
	if status.server.clientExistsByTeamName(game.White) &&
		status.server.clientExistsByTeamName(game.Black) &&
		status.waitCause != allPayersReady {
		return true
	}
	return false
}

// changeCausePlay change cause and server states depending on connecting and disconnecting clients
func (status *status) changeCausePlay() {
	if status.isWaitBothPlayers() {
		status.setWaitBothPlayers()
	} else if status.isWaitWhitePlayer() {
		status.setWaitWhitePlayer()
	} else if status.isWaitBlackPlayer() {
		status.setWaitBlackPlayer()
	} else if status.isAllPlayersReady() {
		status.setAllPlayersReady()
	}
}

// setOverCauseToWhite sets the reason for the end of the game to white win
func (status *status) setOverCauseToWhite() {
	status.overCause = winnerIsWhite
	status.changeOver()
}

// setOverCauseToBlack sets the reason for the end of the game to black win
func (status *status) setOverCauseToBlack() {
	status.overCause = winnerIsBlack
	status.changeOver()
}

// setOverCauseToStalemate sets the reason for ending the game to stalemate
func (status *status) setOverCauseToStalemate() {
	status.overCause = stalemate
	status.changeOver()
}

// setOverCauseToDraw sets game end reason to draw
func (status *status) setOverCauseToDraw() {
	status.overCause = drawGame
	status.changeOver()
}

// changeOver changes the value of the termination state to true if the values of the termination reason contribute to it and sends the data to the broadcast
func (status *status) changeOver() {
	if status.overCause != winnerNotSet {
		status.over = true
	}
	status.send(status.exportOverJSON())
}

// resetOver resets the termination value and termination reason to the default value
func (status *status) resetOver() {
	status.overCause = winnerNotSet
	status.over = false
}

// getCauseOver returns the string value of the termination reason
func (status *status) getCauseOver() string {
	switch status.overCause {
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
func (status *status) exportOverJSON() []byte {
	request := nrp.Simple{Post: "game_over", Body: struct {
		Over  bool   `json:"over"`
		Cause string `json:"cause,omitempty"`
	}{
		Over:  status.over,
		Cause: status.getCauseOver(),
	}}
	return request.Export()
}

// getCausePlay returns the string value of the replay reason
func (status *status) getCausePlay() string {
	switch status.waitCause {
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
func (status *status) exportPlayJSON() []byte {
	request := nrp.Simple{Post: "game_play", Body: struct {
		Play  bool   `json:"play"`
		Cause string `json:"cause,omitempty"`
	}{
		Play:  status.play,
		Cause: status.getCausePlay(),
	}}
	return request.Export()
}

// send data to broadcast
func (status *status) send(exportJSON []byte) {
	status.server.broadcast <- exportJSON
}
