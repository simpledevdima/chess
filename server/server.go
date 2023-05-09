package server

import (
	"fmt"
	"github.com/simpledevdima/chess/game"
	"log"
	"net/http"
)

// newServer returns a server type variable with information for playing chess
func newServer(configFile string) *server {
	return &server{
		config:           newConfig(configFile),
		board:            game.NewBoard(),
		turn:             newTurn(),
		status:           newStatus(),
		timers:           newTimers(),
		drawAttemptsLeft: newDrawAttemptsLeft(),
		clients:          make(map[*client]bool),
		register:         make(chan *client),
		unregister:       make(chan *client),
		broadcast:        make(chan []byte),
	}
}

// server type of server data
type server struct {
	config           *config
	board            *game.Board
	clients          map[*client]bool
	register         chan *client
	unregister       chan *client
	broadcast        chan []byte
	status           *status
	turn             *turn
	timers           *timers
	drawAttemptsLeft *drawAttemptsLeft
}

// newTimes returns a reference to the new timers structure
func newTimers() *timers {
	return &timers{
		white: newTimer(),
		black: newTimer(),
	}
}

// timers struct for timers of team
type timers struct {
	white *timer
	black *timer
}

// newDrawAttemptsLeft returns a reference to the new drawAttemptsLeft structure
func newDrawAttemptsLeft() *drawAttemptsLeft {
	return &drawAttemptsLeft{}
}

// drawAttemptsLeft type with the number of pops to offer a draw
type drawAttemptsLeft struct {
	white int
	black int
}

// setLeft set to drawAttempts time left from argument
func (d *drawAttemptsLeft) setLeft(left int) {
	d.white = left
	d.black = left
}

// loadConfig read the configuration from the installed file and set the data to the structure
func (s *server) loadConfig() {
	s.config.importYAML(s.config.read())
}

// newGame setting server settings for a new game
func (s *server) newGame() {
	// reset over status
	s.status.resetOver()

	// change timer values to values from config
	s.timers.white.setLeft(s.config.StepTimeLeft, s.config.ReserveTimeLeft)
	s.timers.black.setLeft(s.config.StepTimeLeft, s.config.ReserveTimeLeft)

	// set draw attempts left to values from config
	s.drawAttemptsLeft.setLeft(s.config.OfferDrawTimesLeft)

	// new board
	s.board.NewBoard()

	// set turn to default
	s.turn.setDefault()
}

// setupBoardLinks setting links in objects contained inside the board
func (s *server) setLinks() {
	s.status.setServer(s)
	s.timers.white.setTeam(s.board.White)
	s.timers.black.setTeam(s.board.Black)
	s.timers.white.setServer(s)
	s.timers.black.setServer(s)
	s.board.White.SetEnemy(s.board.Black)
	s.board.Black.SetEnemy(s.board.White)
	s.turn.setServer(s)
}

// getFreeTeam return free team as game.TeamName order: white, black, spectators
func (s *server) getFreeTeam() game.TeamName {
	switch {
	case !s.clientExistsByTeamName(game.White):
		return game.White
	case !s.clientExistsByTeamName(game.Black):
		return game.Black
	default:
		return game.Spectators
	}
}

// handlers run server handlers
func (s *server) handlers() {
	// Client registration on server
	http.HandleFunc(s.config.WebsocketURL, func(w http.ResponseWriter, r *http.Request) {
		// making client
		client := &client{server: s, send: make(chan []byte, 256)}
		switch s.getFreeTeam() {
		case game.White:
			// add white team link to client
			client.team = s.board.White
			client.makeDraw()
		case game.Black:
			// add black team link to client
			client.team = s.board.Black
			client.makeDraw()
		default:
			// add spectators team link to client
			client.team = s.board.Spectators
		}
		client.register(w, r)
	})
}

// run listening on a port
func (s *server) run() {
	err := http.ListenAndServe(s.config.Addr, nil)
	if err != nil {
		log.Println(err)
	}
}

// setClientsEnemyLinks setting playing clients links to the opponent's client
func (s *server) setClientsEnemyLinks() {
	wc := s.getClientByTeamName(game.White)
	bc := s.getClientByTeamName(game.Black)
	wc.enemy = bc
	bc.enemy = wc
}

// getClientByTeamName get a reference to the client by specifying the command name in the argument
func (s *server) getClientByTeamName(teamName game.TeamName) *client {
	for client := range s.clients {
		if client.team.Name == teamName {
			return client
		}
	}
	log.Println(fmt.Sprintf("client by %s not found\n", teamName.String()))
	return nil
}

// clientExistsByTeamName returns true if the command has a connected client otherwise returns false
func (s *server) clientExistsByTeamName(teamName game.TeamName) bool {
	for client := range s.clients {
		if client.team.Name == teamName {
			return true
		}
	}
	return false
}

// runChannelProcessing start server
func (s *server) runChannelProcessing() {
	for {
		go s.status.changeCausePlay()
		select {
		case client := <-s.register:
			s.clients[client] = true
			fmt.Println("client connected")
			client.sendGameData()
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				// remove links to the client
				if client.draw != nil {
					client.draw.unsetClient()
				}
				fmt.Println("client disconnected")
			}
		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

// sendGameDataToAll sending game data to all clients
func (s *server) sendGameDataToAll() {
	for client := range s.clients {
		client.sendGameData()
	}
}

// play game on the server
func (s *server) play() {
	if s.turn.now() == game.White {
		go s.timers.white.play()
	} else if s.turn.now() == game.Black {
		go s.timers.black.play()
	}
}

// stop game on the server
func (s *server) stop() {
	if s.turn.now() == game.White {
		s.timers.white.stop()
	} else if s.turn.now() == game.Black {
		s.timers.black.stop()
	}
}

// swapTeams swap clients and their teams
func (s *server) swapTeams() {
	if s.clientExistsByTeamName(game.White) && s.clientExistsByTeamName(game.Black) {
		wc := s.getClientByTeamName(game.White)
		bc := s.getClientByTeamName(game.Black)
		wc, bc = bc, wc
		wc.team, bc.team = bc.team, wc.team
	}
}
