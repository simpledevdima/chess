package chess_bot

import (
	"fmt"
	"github.com/skvdmt/nrp"
)

type status struct {
	teamName  string
	enemyName string
	turn      string
	play      bool
	over      bool
	bot       *Bot
}

// setBot set link to the Bot
func (status *status) setBot(bot *Bot) {
	status.bot = bot
}

// checkForChanges
func (status *status) checkForChanges() {
	if status.isPlay() && !status.isOver() && status.isYourTurn() {
		//fmt.Println("you have to make a move")
		status.bot.move()
	} else {
		//fmt.Println("no move", status.isPlay(), status.isOver(), status.isYourTurn())
	}
}

// setEnemyName
func (status *status) setEnemyName() {
	switch status.teamName {
	case "white":
		status.enemyName = "black"
	case "black":
		status.enemyName = "white"
	}
}

// setTeamName
func (status *status) setTeamName(request *nrp.Simple) {
	var body struct {
		TeamName string `json:"team_name"`
	}
	request.BodyToVariable(&body)
	status.teamName = body.TeamName
	fmt.Println("teamName:", status.teamName)
}

// setTurn
func (status *status) setTurn(request *nrp.Simple) {
	var body struct {
		White bool `json:"white"`
		Black bool `json:"black"`
	}
	request.BodyToVariable(&body)
	if body.White {
		status.turn = "white"
	} else if body.Black {
		status.turn = "black"
	}
	fmt.Println("turn:", status.turn)
}

// setPlay
func (status *status) setPlay(request *nrp.Simple) {
	var body struct {
		Cause string `json:"cause"`
		Play  bool   `json:"play"`
	}
	request.BodyToVariable(&body)
	status.play = body.Play
	fmt.Println("play:", status.play)
}

// setOver
func (status *status) setOver(request *nrp.Simple) {
	var body struct {
		Cause string `json:"cause"`
		Over  bool   `json:"over"`
	}
	request.BodyToVariable(&body)
	status.over = body.Over
	fmt.Println("over:", status.play)
}

// isYourTurn
func (status *status) isYourTurn() bool {
	return status.teamName == status.turn
}

// isPlay
func (status *status) isPlay() bool {
	return status.play
}

// isOver
func (status *status) isOver() bool {
	return status.over
}
