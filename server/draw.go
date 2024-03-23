package server

import (
	"fmt"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
	"log"
	"time"
)

// draw data type with information and data processing methods that allow fixing a draw
type draw struct {
	timeLeftForConfirm int
	open               bool
	client             *client
	requestId          *int
	ticker             *time.Ticker
}

// acceptADraw confirmation for a draw and its installation
func (d *draw) acceptADraw() {
	d.client.enemy.draw.open = false
	d.client.enemy.response(*d.client.enemy.draw.requestId, true, "draw offer accepted")

	// set draw
	d.client.server.stop()
	d.client.server.status.setOverCauseToDraw()
}

// rejectADraw refusal to accept a draw
func (d *draw) rejectADraw() {
	d.client.enemy.draw.open = false
	d.client.enemy.response(*d.client.enemy.draw.requestId, false, "draw offer rejected")
}

// isOpen returns true and an empty string if the draw offer is open otherwise returns a false and a string indicating the reason
func (d *draw) isOpen() (bool, string) {
	return d.open, func() string {
		if !d.open {
			return "draw offer closed"
		}
		return ""
	}()
}

// setRequestId sets a link to the request.Id
func (d *draw) setRequestId(requestId *int) {
	d.requestId = requestId
}

// tick executed after one second has elapsed after receiving a draw offer from the opponent
func (d *draw) tick() {
	d.client.enemy.draw.write(d.exportLeftTimeToConfirmJSON())
	d.timeLeftForConfirm--
	if d.timeLeftForConfirm < 0 {
		// draw time is over
		d.open = false
		d.ticker.Stop()
		d.client.response(*d.requestId, false, "draw offer rejected")
	}
}

// waitResponse countdown for waiting for a response in case of no response at the end of the time, reject the offer
func (d *draw) waitResponse() {
	d.resetTimeLeftForConfirm()
	d.ticker = time.NewTicker(time.Second)
	d.tick()
	for {
		if !d.open {
			break
		}
		select {
		case <-d.ticker.C:
			d.tick()
		}
	}
}

// exportLeftTimeToConfirmJSON returns data on the amount of time left to decide on the confirmation of a draw in JSON format
func (d *draw) exportLeftTimeToConfirmJSON() []byte {
	request := nrp.Simple{Post: "draw_confirm_time", Body: struct {
		Left int `json:"left"`
	}{
		Left: d.timeLeftForConfirm,
	}}
	return request.Export()
}

// resetTimeLeftForConfirm resets the amount of time to make a decision to the value from the configuration
func (d *draw) resetTimeLeftForConfirm() {
	d.timeLeftForConfirm = d.client.server.config.TimeLeftForConfirmDraw
}

// setClient sets the link to the client
func (d *draw) setClient(client *client) {
	d.client = client
}

// unsetClient remove link to the client
func (d *draw) unsetClient() {
	d.client = nil
}

// isValid returns true and an empty string if a draw can be offered otherwise returns false and a string indicating the reason why this is not possible
func (d *draw) isValid() (bool, string) {
	if d.client.server.status.isOver() {
		return false, "game over"
	} else {
		if d.client.server.status.isPlay() {
			if d.open {
				return false, "draw offer already open"
			} else {
				if d.client.server.drawAttemptsLeft.white > 0 && d.client.team.Name == game.White ||
					d.client.server.drawAttemptsLeft.black > 0 && d.client.team.Name == game.Black {
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
func (d *draw) offerADrawToOpponent() {
	d.open = true
	d.client.enemy.draw.write(d.exportOfferADrawToOpponentJSON())
	go d.waitResponse()
	d.reduceAttemptsLeft()
}

// exportOfferADrawToOpponentJSON returns a structure with a request for a draw offer in JSON format
func (d *draw) exportOfferADrawToOpponentJSON() []byte {
	request := nrp.Simple{Post: "opponent_offer_a_draw"}
	return request.Export()
}

// reduceAttemptsLeft reduces the number of attempts to offer a draw
func (d *draw) reduceAttemptsLeft() {
	switch d.client.team.Name {
	case game.White:
		d.client.server.drawAttemptsLeft.white--
	case game.Black:
		d.client.server.drawAttemptsLeft.black--
	}
	d.write(d.exportAttemptsLeftJSON())
}

// write send data to websocket chan of client
func (d *draw) write(data []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovering from d.write()", r)
		}
	}()
	d.client.send <- data
}

// exportAttemptsLeftJSON returns the remaining number of attempts to offer a draw
func (d *draw) exportAttemptsLeftJSON() []byte {
	request := nrp.Simple{Post: "attempts_to_offer_a_draw", Body: struct {
		Left int `json:"left"`
	}{
		Left: func() int {
			switch d.client.team.Name {
			case game.White:
				return d.client.server.drawAttemptsLeft.white
			case game.Black:
				return d.client.server.drawAttemptsLeft.black
			default:
				log.Println("error: unknown team in draw")
				return -1
			}
		}(),
	}}
	return request.Export()
}
