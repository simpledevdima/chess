package chess_bot

import (
	"fmt"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
)

func NewMove(bot *Bot, from *game.Position, to *game.Position) *move {
	m := new(move)
	m.setBot(bot)
	if from != nil {
		m.From.Set(from.X, from.Y)
	}
	if to != nil {
		m.To.Set(to.X, to.Y)
	}
	return m
}

type move struct {
	From struct {
		game.Position `json:"position"`
	} `json:"from"`
	To struct {
		game.Position `json:"position"`
	} `json:"to"`
	bot *Bot
}

// setBot
func (move *move) setBot(bot *Bot) {
	move.bot = bot
}

// send отправка данных о совершаемом ходе на сервер
func (move *move) send() {
	if move.bot.status.isYourTurn() {
		request := nrp.Simple{Post: "move", Body: move}
		fmt.Println("BOT SEND MOVE:", string(request.Export()))
		move.bot.send <- request.Export()
	}
}

// exec выполнение хода информаци о котором поступила с сервера
func (move *move) exec() {
	//fmt.Println(move)
	if move.bot.team.Figures.ExistsByCoords(move.From.Position.X, move.From.Position.Y) {
		// your team move
		figureID := move.bot.team.Figures.GetIndexByCoords(move.From.Position.X, move.From.Position.Y)
		move.bot.team.Figures[figureID].Move(move.To.Position.X, move.To.Position.Y)
		//fmt.Println(move.bot.team.Figures)
	} else if move.bot.enemy.Figures.ExistsByCoords(move.From.Position.X, move.From.Position.Y) {
		// enemy team move
		figureID := move.bot.enemy.Figures.GetIndexByCoords(move.From.Position.X, move.From.Position.Y)
		move.bot.enemy.Figures[figureID].Move(move.To.Position.X, move.To.Position.Y)
		//fmt.Println(move.bot.enemy.Figures)
	}
}

//if move.bot.data.status.turn == move.bot.data.status.teamName {
//
//} else {
//
//}
//figureID, err := move.bot.team.GetFigureID(move.From.Position.X, move.From.Position.Y)
//if err != nil {
//	log.Println(err)
//}
//
//move.bot.team.Figures[figureID].SetPosition(move.To.Position.X, move.To.Position.Y)
//if move.bot.enemy.FigureExist(move.To.Position.X, move.To.Position.Y) {
//	// eat enemy figure
//	err := move.bot.enemy.Eating(move.To.Position.X, move.To.Position.Y)
//	if err != nil {
//		log.Println(err)
//	}
//}
