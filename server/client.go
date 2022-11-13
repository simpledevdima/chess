package server

import (
	"github.com/gorilla/websocket"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
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

// makeDraw create a draw in the client and set draw to a reference to the client
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
	request := nrp.Simple{Post: "your", Body: struct {
		TeamName string `json:"team_name"`
	}{
		TeamName: client.team.Name.String(),
	}}
	return request.Export()
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

// response send response to request to client
func (client *client) response(id int, valid bool, cause string) {
	response := &nrp.Simple{Post: "response", Body: &struct {
		RequestId int    `json:"request_id,omitempty"`
		Valid     bool   `json:"valid"`
		Cause     string `json:"cause,omitempty"`
	}{
		RequestId: id,
		Valid:     valid,
		Cause:     cause,
	}}
	client.send <- response.Export()
}

// postToMove if a request for a move came
func (client *client) postToMove(request *nrp.Simple) {
	var move move
	move.setClient(client)
	request.BodyToVariable(&move)
	valid, cause := move.isValid()
	client.response(request.Id, valid, cause)
	if valid {
		move.exec()
		client.server.turn.change()
	}
}

// postToSurrender if the request came to surrender
func (client *client) postToSurrender(request *nrp.Simple) {
	var surrender surrender
	surrender.setClient(client)
	valid, cause := surrender.isValid()
	if valid {
		surrender.exec()
	}
	client.response(request.Id, valid, cause)
}

// postToNewGame if you are asked to create a new game
func (client *client) postToNewGame(request *nrp.Simple) {
	var newGame newGame
	newGame.setServer(client.server)
	valid, cause := newGame.isValid()
	if valid {
		newGame.exec()
	}
	client.response(request.Id, valid, cause)
}

// postToOfferADraw if a request came with a draw offer
func (client *client) postToOfferADraw(request *nrp.Simple) {
	valid, cause := client.draw.isValid()
	if valid {
		client.draw.setRequestId(&request.Id)
		client.draw.offerADrawToOpponent()
	} else {
		client.response(request.Id, valid, cause)
	}
}

// postToDrawOfferAccepted if a request is received approving the offer of a draw
func (client *client) postToDrawOfferAccepted(request *nrp.Simple) {
	valid, cause := client.enemy.draw.isOpen()
	if valid {
		client.draw.acceptADraw()
	}
	client.response(request.Id, valid, cause)
}

// postToDrawOfferRejected if a request is received rejecting the offer of a draw
func (client *client) postToDrawOfferRejected(request *nrp.Simple) {
	valid, cause := client.enemy.draw.isOpen()
	if valid {
		client.draw.rejectADraw()
	}
	client.response(request.Id, valid, cause)
}

// incomingDataProcessing handles the request from the argument
func (client *client) incomingDataProcessing(dataJSON []byte) {
	var request nrp.Simple
	request.Parse(dataJSON)
	switch client.team.Name {
	case game.White, game.Black:
		switch request.Post {
		case "move":
			client.postToMove(&request)
		case "surrender":
			client.postToSurrender(&request)
		case "new":
			client.postToNewGame(&request)
		case "offer_a_draw":
			client.postToOfferADraw(&request)
		case "draw_offer_accepted":
			client.postToDrawOfferAccepted(&request)
		case "draw_offer_rejected":
			client.postToDrawOfferRejected(&request)
		default:
			client.response(request.Id, false, "unknown request")
		}
	default:
		client.response(request.Id, false, "you are a spectator and cannot send requests")
	}
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
		_, dataJSDON, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		client.incomingDataProcessing(dataJSDON)
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
	if client.server.status.isPlay() && !client.server.status.isOver() {
		switch client.server.turn.now() {
		case game.White:
			client.server.timers.white.stop()
		case game.Black:
			client.server.timers.black.stop()
		}
	}
	switch client.team.Name {
	case game.White:
		client.server.status.setOverCauseToBlack()
	case game.Black:
		client.server.status.setOverCauseToWhite()
	}
}
