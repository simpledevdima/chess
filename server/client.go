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
func (c *client) makeDraw() {
	c.draw = &draw{}
	c.draw.setClient(c)
}

// sendGameData sends game data to the client
func (c *client) sendGameData() {
	c.send <- c.exportJSON()
	c.send <- c.server.board.ExportJSON()
	c.send <- c.server.turn.exportJSON()
	c.send <- c.server.timers.white.exportReserveJSON()
	c.send <- c.server.timers.black.exportReserveJSON()
	switch c.server.turn.now() {
	case game.White:
		c.send <- c.server.timers.white.exportStepJSON()
	case game.Black:
		c.send <- c.server.timers.black.exportStepJSON()
	}
	switch c.team.Name {
	case game.White, game.Black:
		c.send <- c.draw.exportAttemptsLeftJSON()
	}
	c.send <- c.server.status.exportPlayJSON()
	c.send <- c.server.status.exportOverJSON()
}

// exportJSON getting data with the name of the command of the current client in JSON format
func (c *client) exportJSON() []byte {
	request := nrp.Simple{Post: "your", Body: struct {
		TeamName string `json:"team_name"`
	}{
		TeamName: c.team.Name.String(),
	}}
	return request.Export()
}

// register making websocket connection with client, run read and write methods and register current client on server
func (c *client) register(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			var origin = r.Header.Get("origin")
			if origin == c.server.config.OriginalClientURL {
				return true
			}
			return false
		},
	}
	var err error
	c.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	c.server.register <- c

	go c.read()
	go c.write()
}

// response send response to request to client
func (c *client) response(id int, valid bool, cause string) {
	response := &nrp.Simple{Post: "response", Body: &struct {
		RequestId int    `json:"request_id,omitempty"`
		Valid     bool   `json:"valid"`
		Cause     string `json:"cause,omitempty"`
	}{
		RequestId: id,
		Valid:     valid,
		Cause:     cause,
	}}
	c.send <- response.Export()
}

// postToMove if a request for a move came
func (c *client) postToMove(request *nrp.Simple) {
	m := newMove(c)
	request.BodyToVariable(&m)
	valid, cause := m.isValid()
	c.response(request.Id, valid, cause)
	if valid {
		m.exec()
		c.server.turn.change()
	}
}

// postToSurrender if the request came to surrender
func (c *client) postToSurrender(request *nrp.Simple) {
	var surrender surrender
	surrender.setClient(c)
	valid, cause := surrender.isValid()
	if valid {
		surrender.exec()
	}
	c.response(request.Id, valid, cause)
}

// postToNewGame if you are asked to create a new game
func (c *client) postToNewGame(request *nrp.Simple) {
	var newGame newGame
	newGame.setServer(c.server)
	valid, cause := newGame.isValid()
	if valid {
		newGame.exec()
	}
	c.response(request.Id, valid, cause)
}

// postToOfferADraw if a request came with a draw offer
func (c *client) postToOfferADraw(request *nrp.Simple) {
	valid, cause := c.draw.isValid()
	if valid {
		c.draw.setRequestId(&request.Id)
		c.draw.offerADrawToOpponent()
	} else {
		c.response(request.Id, valid, cause)
	}
}

// postToDrawOfferAccepted if a request is received approving the offer of a draw
func (c *client) postToDrawOfferAccepted(request *nrp.Simple) {
	valid, cause := c.enemy.draw.isOpen()
	if valid {
		c.draw.acceptADraw()
	}
	c.response(request.Id, valid, cause)
}

// postToDrawOfferRejected if a request is received rejecting the offer of a draw
func (c *client) postToDrawOfferRejected(request *nrp.Simple) {
	valid, cause := c.enemy.draw.isOpen()
	if valid {
		c.draw.rejectADraw()
	}
	c.response(request.Id, valid, cause)
}

// incomingDataProcessing handles the request from the argument
func (c *client) incomingDataProcessing(dataJSON []byte) {
	var request nrp.Simple
	request.Parse(dataJSON)
	switch c.team.Name {
	case game.White, game.Black:
		switch request.Post {
		case "move":
			c.postToMove(&request)
		case "surrender":
			c.postToSurrender(&request)
		case "new":
			c.postToNewGame(&request)
		case "offer_a_draw":
			c.postToOfferADraw(&request)
		case "draw_offer_accepted":
			c.postToDrawOfferAccepted(&request)
		case "draw_offer_rejected":
			c.postToDrawOfferRejected(&request)
		default:
			c.response(request.Id, false, "unknown request")
		}
	default:
		c.response(request.Id, false, "you are a spectator and cannot send requests")
	}
}

// read receive request data from client and execute requests
func (c *client) read() {
	defer func() {
		c.server.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		_, dataJSDON, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		c.incomingDataProcessing(dataJSDON)
	}
}

// write send data from client channel to websocket connection
func (c *client) write() {
	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// surrender the client accepts defeat and gives victory to the enemy
func (c *client) surrender() {
	if c.server.status.isPlay() && !c.server.status.isOver() {
		switch c.server.turn.now() {
		case game.White:
			c.server.timers.white.stop()
		case game.Black:
			c.server.timers.black.stop()
		}
	}
	switch c.team.Name {
	case game.White:
		c.server.status.setOverCauseToBlack()
	case game.Black:
		c.server.status.setOverCauseToWhite()
	}
}
