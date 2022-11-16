package server

import (
	"github.com/skvdmt/chess/game"
)

// surrender surrender control capability data type
type surrender struct {
	client *client
}

// isValid returns true and an empty string if you can give up otherwise returns false and a reason why you can't give up
func (s *surrender) isValid() (bool, string) {
	if s.client.server.status.isOver() {
		return false, "game is over"
	} else {
		switch s.client.team.Name {
		case game.White, game.Black:
			return true, ""
		default:
			return false, "you cant surrender"
		}
	}
}

// setClient setting a client link
func (s *surrender) setClient(client *client) {
	s.client = client
}

// exec delivery
func (s *surrender) exec() {
	s.client.surrender()
}
