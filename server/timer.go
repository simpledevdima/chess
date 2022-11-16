package server

import (
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
	"time"
)

func newTimer() *timer {
	return &timer{}
}

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
func (t *timer) setLeft(stepLeft int, reserveLeft int) {
	t.stepLeft = stepLeft
	t.reserveLeft = reserveLeft
}

// setTeam set link to team
func (t *timer) setTeam(team *game.Team) {
	t.team = team
}

// setServer set link to server
func (t *timer) setServer(server *server) {
	t.server = server
}

// tick executed after passing one second of the turn
func (t *timer) tick() {
	if t.stepLeft >= 0 {
		t.send(t.exportStepJSON())
		t.stepLeft--
	} else {
		if t.server.status.isPlay() {
			t.reserveLeft--
			t.send(t.exportReserveJSON())
			if t.reserveLeft == 0 {
				// time over
				t.stop()
				switch t.team.Name {
				case game.White:
					t.server.status.setOverCauseToBlack()
				case game.Black:
					t.server.status.setOverCauseToWhite()
				}
			}
		}
	}
}

// play start of calculation of command turn time
func (t *timer) play() {
	t.stopFlag = false
	t.ticker = time.NewTicker(time.Second)
	t.tick()
	for {
		select {
		case <-t.ticker.C:
			t.tick()
		}
		if t.stopFlag {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// isOver return true if reserve time of team is over
func (t *timer) isOver() bool {
	if t.reserveLeft > 0 {
		return false
	}
	return true
}

// stop set stopFlag to true for break countdown
func (t *timer) stop() {
	t.ticker.Stop()
	t.stopFlag = true
}

// reset stepTimeLeft to count from config
func (t *timer) reset() {
	t.stepLeft = t.server.config.StepTimeLeft
}

// exportStepJSON return stepTimeLeft data from struct in JSON
func (t *timer) exportStepJSON() []byte {
	request := nrp.Simple{Post: "step_time", Body: struct {
		TeamName string `json:"team_name"`
		Left     int    `json:"left"`
	}{
		TeamName: t.team.Name.String(),
		Left:     t.getTimeLeft(),
	}}
	return request.Export()
}

// exportReserveJSON return reserveTimeLeft data from struct in JSON
func (t *timer) exportReserveJSON() []byte {
	request := nrp.Simple{Post: "reserve_time", Body: struct {
		TeamName string `json:"team_name"`
		Left     int    `json:"left"`
	}{
		TeamName: t.team.Name.String(),
		Left:     t.reserveLeft,
	}}
	return request.Export()
}

// getTimLeft return time left
func (t *timer) getTimeLeft() int {
	if t.stepLeft >= 0 {
		return t.stepLeft
	} else {
		return t.reserveLeft
	}
}

// send data to broadcast
func (t *timer) send(dataJSON []byte) {
	t.server.broadcast <- dataJSON
}
