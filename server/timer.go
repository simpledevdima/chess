package server

import (
	"encoding/json"
	"github.com/skvdmt/chess/game"
	"log"
	"time"
)

// timer data structure containing the remaining time of the command
type timer struct {
	stepLeft    int
	reserveLeft int
	stopFlag    bool
	team        *game.Team
	server      *server
	ticker      *time.Ticker
}

// setLeft set the remaining time of the command
func (timer *timer) setLeft(stepLeft int, reserveLeft int) {
	timer.stepLeft = stepLeft
	timer.reserveLeft = reserveLeft
}

// setTeam set link to team
func (timer *timer) setTeam(team *game.Team) {
	timer.team = team
}

// setServer set link to server
func (timer *timer) setServer(server *server) {
	timer.server = server
}

// tick executed after passing one second of the turn
func (timer *timer) tick() {
	if timer.stepLeft >= 0 {
		timer.send(timer.exportStepJSON())
		timer.stepLeft--
	} else {
		if timer.server.status.isPlay() {
			timer.reserveLeft--
			timer.send(timer.exportReserveJSON())
			if timer.reserveLeft == 0 {
				// time over
				timer.stop()
				switch timer.team.Name {
				case game.White:
					timer.server.status.setOverCauseToBlack()
				case game.Black:
					timer.server.status.setOverCauseToWhite()
				}
			}
		}
	}
}

// play start of calculation of command turn time
func (timer *timer) play() {
	timer.stopFlag = false
	timer.ticker = time.NewTicker(time.Second)
	timer.tick()
	for {
		select {
		case <-timer.ticker.C:
			timer.tick()
		}
		if timer.stopFlag {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// isOver return true if reserve time of team is over
func (timer *timer) isOver() bool {
	if timer.reserveLeft > 0 {
		return false
	}
	return true
}

// stop set stopFlag to true for break countdown
func (timer *timer) stop() {
	timer.ticker.Stop()
	timer.stopFlag = true
}

// reset stepTimeLeft to count from config
func (timer *timer) reset() {
	timer.stepLeft = timer.server.config.StepTimeLeft
}

// exportStepJSON return stepTimeLeft data from struct in JSON
func (timer *timer) exportStepJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		TeamName     string `json:"team_name"`
		StepTimeLeft int    `json:"step_time_left"`
	}{
		TeamName:     timer.team.Name.String(),
		StepTimeLeft: timer.getTimeLeft(),
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// exportReserveJSON return reserveTimeLeft data from struct in JSON
func (timer *timer) exportReserveJSON() []byte {
	dataJSON, err := json.Marshal(struct {
		TeamName        string `json:"team_name"`
		ReserveTimeLeft int    `json:"reserve_time_left"`
	}{
		TeamName:        timer.team.Name.String(),
		ReserveTimeLeft: timer.reserveLeft,
	})
	if err != nil {
		log.Println(err)
	}
	return dataJSON
}

// getTimLeft return time left
func (timer *timer) getTimeLeft() int {
	if timer.stepLeft >= 0 {
		return timer.stepLeft
	} else {
		return timer.reserveLeft
	}
}

// send data to broadcast
func (timer *timer) send(dataJSON []byte) {
	timer.server.broadcast <- dataJSON
}
