package server

import (
	"encoding/json"
	"fmt"
	"github.com/skvdmt/chess/game"
	"log"
	"time"
)

// draw data type with information and data processing methods that allow fixing a draw
type draw struct {
	timeLeftForConfirm int
	open               bool
	client             *client
	request            *request
	ticker             *time.Ticker
}

// acceptADraw confirmation for a draw and its installation
func (draw *draw) acceptADraw() {
	draw.client.enemy.draw.open = false
	draw.client.enemy.draw.request.makeAndSendResponse(draw.client.enemy.draw.request.getResponseValid(true, "draw offer accepted"))

	// set draw
	draw.client.server.stop()
	draw.client.server.status.setOverCauseToDraw()
}

// rejectADraw refusal to accept a draw
func (draw *draw) rejectADraw() {
	draw.client.enemy.draw.open = false
	draw.client.enemy.draw.request.makeAndSendResponse(draw.client.enemy.draw.request.getResponseValid(false, "draw offer rejected"))
}

// isOpen returns true and an empty string if the draw offer is open otherwise returns a false and a string indicating the reason
func (draw *draw) isOpen() (bool, string) {
	return draw.open, func() string {
		if !draw.open {
			return "draw offer closed"
		}
		return ""
	}()
}

// setRequest sets a link to the request
func (draw *draw) setRequest(request *request) {
	draw.request = request
}

// tick executed after one second has elapsed after receiving a draw offer from the opponent
func (draw *draw) tick() {
	draw.client.enemy.draw.write(draw.exportLeftTimeToConfirmJSON())
	draw.timeLeftForConfirm--
	if draw.timeLeftForConfirm < 0 {
		// draw time is over
		draw.open = false
		draw.ticker.Stop()
		draw.request.makeAndSendResponse(draw.request.getResponseValid(false, "draw offer rejected"))
	}
}

// waitResponse countdown for waiting for a response in case of no response at the end of the time, reject the offer
func (draw *draw) waitResponse() {
	draw.resetTimeLeftForConfirm()
	draw.ticker = time.NewTicker(time.Second)
	draw.tick()
	for {
		if !draw.open {
			break
		}
		select {
		case <-draw.ticker.C:
			draw.tick()
		}
	}
}

// exportLeftTimeToConfirmJSON returns data on the amount of time left to decide on the confirmation of a draw in JSON format
func (draw *draw) exportLeftTimeToConfirmJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		TimeLeftForConfirmDraw int `json:"time_left_for_confirm_draw"`
	}{
		TimeLeftForConfirmDraw: draw.timeLeftForConfirm,
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// resetTimeLeftForConfirm resets the amount of time to make a decision to the value from the configuration
func (draw *draw) resetTimeLeftForConfirm() {
	draw.timeLeftForConfirm = draw.client.server.config.TimeLeftForConfirmDraw
}

// setClient sets the link to the client
func (draw *draw) setClient(client *client) {
	draw.client = client
}

// isValid returns true and an empty string if a draw can be offered otherwise returns false and a string indicating the reason why this is not possible
func (draw *draw) isValid() (bool, string) {
	if draw.client.server.status.isOver() {
		return false, "game over"
	} else {
		if draw.client.server.status.isPlay() {
			if draw.open {
				return false, "draw offer already open"
			} else {
				if draw.client.server.drawAttemptsLeft.white > 0 && draw.client.team.Name == game.White ||
					draw.client.server.drawAttemptsLeft.black > 0 && draw.client.team.Name == game.Black {
					return true, ""
				}
				return false, "attempts to offer a draw ended"
			}
		} else {
			return false, "game are stopped"
		}
	}
}

// offerADrawToOpponent is executed when a request for a draw is received from the opponent
func (draw *draw) offerADrawToOpponent() {
	draw.open = true
	draw.client.enemy.draw.write(draw.exportOfferADrawToOpponentJSON())
	go draw.waitResponse()
	draw.reduceAttemptsLeft()
}

// exportOfferADrawToOpponentJSON returns a structure with a request for a draw offer in JSON format
func (draw *draw) exportOfferADrawToOpponentJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		OpponentOfferADraw bool `json:"opponent_offer_a_draw"`
	}{
		OpponentOfferADraw: true,
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// reduceAttemptsLeft reduces the number of attempts to offer a draw
func (draw *draw) reduceAttemptsLeft() {
	switch draw.client.team.Name {
	case game.White:
		draw.client.server.drawAttemptsLeft.white--
	case game.Black:
		draw.client.server.drawAttemptsLeft.black--
	}
	draw.write(draw.exportAttemptsLeftJSON())
}

// write send data to websocket chan of client
func (draw *draw) write(data []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovering from draw.write()", r)
		}
	}()
	draw.client.send <- data
}

// exportAttemptsLeftJSON returns the remaining number of attempts to offer a draw
func (draw *draw) exportAttemptsLeftJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		AttemptsLeftToOfferADraw int `json:"attempts_left_to_offer_a_draw"`
	}{
		AttemptsLeftToOfferADraw: func() int {
			switch draw.client.team.Name {
			case game.White:
				return draw.client.server.drawAttemptsLeft.white
			case game.Black:
				return draw.client.server.drawAttemptsLeft.black
			default:
				log.Println("error: unknown team in draw")
				return -1
			}
		}(),
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}
