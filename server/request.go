package server

import (
	"encoding/json"
	"github.com/skvdmt/chess/game"
	"log"
)

// request data type of the request to the server and its processing
type request struct {
	id     int
	post   string
	body   []byte
	client *client
}

// setClient set a link to the client in the request
func (request *request) setClient(client *client) {
	request.client = client
}

// importJSON get data from json and set it to request structure
func (request *request) importJSON(jsonData []byte) {
	var iJSON struct {
		Id   int         `json:"id"`
		Post string      `json:"post"`
		Body interface{} `json:"body"`
	}
	err := json.Unmarshal(jsonData, &iJSON)
	if err != nil {
		log.Println(err)
	}
	request.id = iJSON.Id
	request.post = iJSON.Post
	request.body, err = json.Marshal(iJSON.Body)
	if err != nil {
		log.Println(err)
	}
}

// makeAndSendResponse creates a response to the request, puts an argument in its body, and sends the response
func (request *request) makeAndSendResponse(body []byte) {
	var response response
	response.setClient(request.client)
	response.setID(request.id)
	response.setBody(body)
	response.write(response.exportJSON())
}

// getResponseValid returns a response data structure with validity parameters and the reason for this state in JSON format
func (request *request) getResponseValid(valid bool, cause string) []byte {
	var responseValid responseValid
	if valid {
		responseValid.setValid(true)
	} else {
		responseValid.setValid(false)
		responseValid.setCause(cause)
	}
	return responseValid.exportJSON()
}

// exec the client's request and sends a response to it
func (request *request) exec() {
	switch request.client.team.Name {
	case game.White, game.Black:
		switch request.post {
		case "move":
			var move move
			move.setClient(request.client)
			move.importJSON(request.body)
			valid, cause := move.isValid()
			if valid {
				move.exec()
				request.client.server.turn.change()
			}
			request.makeAndSendResponse(request.getResponseValid(valid, cause))
		case "surrender":
			var surrender surrender
			surrender.setClient(request.client)
			valid, cause := surrender.isValid()
			if valid {
				surrender.exec()
			}
			request.makeAndSendResponse(request.getResponseValid(valid, cause))
		case "new":
			var newGame newGame
			newGame.setServer(request.client.server)
			valid, cause := newGame.isValid()
			if valid {
				newGame.exec()
			}
			request.makeAndSendResponse(request.getResponseValid(valid, cause))
		case "offer_a_draw":
			valid, cause := request.client.draw.isValid()
			if valid {
				request.client.draw.setRequest(request)
				request.client.draw.offerADrawToOpponent()
			} else {
				request.makeAndSendResponse(request.getResponseValid(valid, cause))
			}
		case "draw_offer_accepted":
			valid, cause := request.client.enemy.draw.isOpen()
			if valid {
				request.client.draw.acceptADraw()
			}
			request.makeAndSendResponse(request.getResponseValid(valid, cause))
		case "draw_offer_rejected":
			valid, cause := request.client.enemy.draw.isOpen()
			if valid {
				request.client.draw.rejectADraw()
			}
			request.makeAndSendResponse(request.getResponseValid(valid, cause))
		default:
			request.makeAndSendResponse(request.getResponseValid(false, "unknown request"))
		}
	default:
		request.makeAndSendResponse(request.getResponseValid(false, "you are a spectator and cannot send requests"))
	}
}
