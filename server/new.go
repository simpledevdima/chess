package server

// newGame data type for making a new game
type newGame struct {
	server *server
}

// setServer set link to server
func (newGame *newGame) setServer(server *server) {
	newGame.server = server
}

// isValid returns true if it is possible to start a new game and an empty string otherwise returns false and a string with the reason why it is not possible to start a new game
func (newGame *newGame) isValid() (bool, string) {
	if newGame.server.status.isOver() {
		return true, ""
	} else {
		return false, "game not over"
	}
}

// exec making a new game
func (newGame *newGame) exec() {
	if newGame.server.config.SwapTeamsAfterMakingNewGame {
		newGame.server.swapTeams()
	}
	newGame.server.newGame()
	newGame.server.sendGameDataToAll()
	newGame.server.status.changePlay()
}
