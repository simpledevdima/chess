package server

// newGame data type for making a new game
type newGame struct {
	server *server
}

// setServer set link to server
func (n *newGame) setServer(server *server) {
	n.server = server
}

// isValid returns true if it is possible to start a new game and an empty string otherwise returns false and a string with the reason why it is not possible to start a new game
func (n *newGame) isValid() (bool, string) {
	if n.server.status.isOver() {
		return true, ""
	} else {
		return false, "game not over"
	}
}

// exec making a new game
func (n *newGame) exec() {
	if n.server.config.SwapTeamsAfterMakingNewGame {
		n.server.swapTeams()
	}
	n.server.newGame()
	n.server.sendGameDataToAll()
	n.server.status.changePlay()
}
