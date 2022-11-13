// Package server contains a functional chess server that allows two players to play chess, as well as watching it for many spectators
package server

import (
	"github.com/skvdmt/chess/game"
)

// Start API to run chess server
func Start(configFile string) {
	// make server
	server := getServer()

	// get config
	server.config.read(configFile)

	// setup
	server.setLinks()
	server.newGame()

	// run channel processing
	go server.runChannelProcessing()

	// run handlers
	server.handlers()

	// run server
	server.run()
}

// getServer returns a structure of the server type with data for the game process
func getServer() *server {
	return &server{
		config: &config{},
		board: &game.Board{
			White:      &game.Team{Name: game.White},
			Black:      &game.Team{Name: game.Black},
			Spectators: &game.Team{Name: game.Spectators},
		},
		clients:          make(map[*client]bool),
		register:         make(chan *client),
		unregister:       make(chan *client),
		broadcast:        make(chan []byte),
		turn:             &turn{},
		status:           &status{},
		timers:           &timers{white: &timer{}, black: &timer{}},
		drawAttemptsLeft: &drawAttemptsLeft{},
	}
}
