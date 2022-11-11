package server

import (
	"fmt"
	"github.com/skvdmt/chess/game"
	"log"
	"net/http"
)

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

// timers struct for timers of team
type timers struct {
	white *timer
	black *timer
}

// attemptsLeft
type drawAttemptsLeft struct {
	white int
	black int
}

// newGame
func (server *server) newGame() {
	// over
	server.status.resetOver()

	// timers
	server.timers.white.setLeft(server.config.StepTimeLeft, server.config.ReserveTimeLeft)
	server.timers.black.setLeft(server.config.StepTimeLeft, server.config.ReserveTimeLeft)

	// turn
	server.turn.setDefault()

	// draw attempts left
	server.drawAttemptsLeft.white = server.config.OfferDrawTimesLeft
	server.drawAttemptsLeft.black = server.config.OfferDrawTimesLeft

	// new board
	server.board.NewBoard()
}

// setupBoardLinks setting links in objects contained inside the board
func (server *server) setLinks() {
	server.status.setServer(server)
	server.timers.white.setTeam(server.board.White)
	server.timers.black.setTeam(server.board.Black)
	server.timers.white.setServer(server)
	server.timers.black.setServer(server)
	server.board.White.SetEnemy(server.board.Black)
	server.board.Black.SetEnemy(server.board.White)
	server.turn.setServer(server)
}

// getFreeTeam return free team as game.TeamName order: white, black, spectators
func (server *server) getFreeTeam() game.TeamName {
	switch {
	case !server.clientExistsByTeamName(game.White):
		return game.White
	case !server.clientExistsByTeamName(game.Black):
		return game.Black
	default:
		return game.Spectators
	}
}

// handlers run server handlers
func (server *server) handlers() {
	// Client registration on server
	http.HandleFunc(server.config.WebsocketURL, func(w http.ResponseWriter, r *http.Request) {
		// making client
		client := &client{server: server, send: make(chan []byte, 256)}
		switch server.getFreeTeam() {
		case game.White:
			// add white team link to client
			client.team = server.board.White
			client.makeDraw()
		case game.Black:
			// add black team link to client
			client.team = server.board.Black
			client.makeDraw()
		default:
			// add spectators team link to client
			client.team = server.board.Spectators
		}
		client.register(w, r)
	})
}

// run listening on a port
func (server *server) run() {
	err := http.ListenAndServe(server.config.Addr, nil)
	if err != nil {
		log.Println(err)
	}
}

func (server *server) setClientsEnemyLinks() {
	wc := server.getClientByTeamName(game.White)
	bc := server.getClientByTeamName(game.Black)
	wc.enemy = bc
	bc.enemy = wc
}

func (server *server) getClientByTeamName(teamName game.TeamName) *client {
	for client := range server.clients {
		if client.team.Name == teamName {
			return client
		}
	}
	log.Println(fmt.Sprintf("client by %s not found\n", teamName.String()))
	return nil
}

// clientExistsByTeamName returns true if the command has a connected client otherwise returns false
func (server *server) clientExistsByTeamName(teamName game.TeamName) bool {
	for client := range server.clients {
		if client.team.Name == teamName {
			return true
		}
	}
	return false
}

// runChannelProcessing start server
func (server *server) runChannelProcessing() {
	for {
		go server.status.changeCausePlay()
		select {
		case client := <-server.register:
			server.clients[client] = true
			fmt.Println("client connected")
			client.sendGameData()
		case client := <-server.unregister:
			if _, ok := server.clients[client]; ok {
				delete(server.clients, client)
				fmt.Println("client disconnected")
			}
		case message := <-server.broadcast:
			for client := range server.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(server.clients, client)
				}
			}
		}
	}
}

// sendGameDataToAll
func (server *server) sendGameDataToAll() {
	for client := range server.clients {
		client.sendGameData()
	}
}

// play game on the server
func (server *server) play() {
	if server.turn.now() == game.White {
		go server.timers.white.play()
	} else if server.turn.now() == game.Black {
		go server.timers.black.play()
	}
}

// stop game on the server
func (server *server) stop() {
	if server.turn.now() == game.White {
		server.timers.white.stop()
	} else if server.turn.now() == game.Black {
		server.timers.black.stop()
	}
}

// swapTeams swap clients and their teams
func (server *server) swapTeams() {
	if server.clientExistsByTeamName(game.White) && server.clientExistsByTeamName(game.Black) {
		wc := server.getClientByTeamName(game.White)
		bc := server.getClientByTeamName(game.Black)
		wc, bc = bc, wc
		wc.team, bc.team = bc.team, wc.team
	}
}
