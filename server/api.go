// Package server contains a functional chess server that allows two players to play chess, as well as watching it for many spectators
package server

// Start API to run chess server
func Start(configFile string) {
	server := newServer(configFile)  // making a new empty server
	server.loadConfig()              // get config from installed configuration file
	server.setLinks()                // setting links inside server types
	server.newGame()                 // setup of new game
	go server.runChannelProcessing() // run channel processing
	server.handlers()                // run handlers
	server.run()                     // run server
}
