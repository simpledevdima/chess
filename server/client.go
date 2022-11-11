package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/skvdmt/chess/game"
	"log"
	"net/http"
)

// client type of data for current client on server
type client struct {
	server *server
	conn   *websocket.Conn
	send   chan []byte
	status string
	team   *game.Team
	enemy  *client
	draw   *draw
}

func (client *client) makeDraw() {
	client.draw = &draw{}
	client.draw.setClient(client)
}

// sendGameData sends game data to the client
func (client *client) sendGameData() {
	client.send <- client.exportJSON()
	client.send <- client.server.board.ExportJSON()
	client.send <- client.server.turn.exportJSON()
	client.send <- client.server.timers.white.exportReserveJSON()
	client.send <- client.server.timers.black.exportReserveJSON()
	switch client.server.turn.now() {
	case game.White:
		client.send <- client.server.timers.white.exportStepJSON()
	case game.Black:
		client.send <- client.server.timers.black.exportStepJSON()
	}
	switch client.team.Name {
	case game.White, game.Black:
		client.send <- client.draw.exportAttemptsLeftJSON()
	}
	client.send <- client.server.status.exportPlayJSON()
	client.send <- client.server.status.exportOverJSON()
}

// exportJSON getting data with the name of the command of the current client in JSON format
func (client *client) exportJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		YourTeamName string `json:"your_team_name"`
	}{
		YourTeamName: client.team.Name.String(),
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// register making websocket connection with client, run read and write methods and register current client on server
func (client *client) register(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			var origin = r.Header.Get("origin")
			if origin == client.server.config.OriginalClientURL {
				return true
			}
			return false
		},
	}
	var err error
	client.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	client.server.register <- client

	go client.read()
	go client.write()
}

// read receive request data from client and execute requests
func (client *client) read() {
	defer func() {
		client.server.unregister <- client
		err := client.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		_, jsonData, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var request request
		request.importJSON(jsonData)
		request.setClient(client)
		request.exec()
	}
}

// write send data from client channel to websocket connection
func (client *client) write() {
	for {
		select {
		case message := <-client.send:
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// surrender the client accepts defeat and gives victory to the enemy
func (client *client) surrender() {
	switch client.server.turn.now() {
	case game.White:
		client.server.timers.white.stop()
	case game.Black:
		client.server.timers.black.stop()
	}
	switch client.enemy.team.Name {
	case game.White:
		client.server.status.setOverCauseToWhite()
	case game.Black:
		client.server.status.setOverCauseToBlack()
	}
}
