package chess_bot

import (
	"fmt"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
)

func newMove(bot *Bot, from *game.Position, to *game.Position) *move {
	m := new(move)
	m.setBot(bot)
	m.From.Position = from
	m.To.Position = to
	return m
}

type move struct {
	From struct {
		*game.Position `json:"position"`
	} `json:"from"`
	To struct {
		*game.Position `json:"position"`
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
	if move.bot.team.Figures.ExistsByPosition(move.From.Position) {
		// your team move
		figureID := move.bot.team.Figures.GetIndexByPosition(move.From.Position)
		move.bot.team.Figures[figureID].Move(move.To.Position)
		//fmt.Println(move.bot.team.Figures)
	} else if move.bot.enemy.Figures.ExistsByPosition(move.From.Position) {
		// enemy team move
		figureID := move.bot.enemy.Figures.GetIndexByPosition(move.From.Position)
		move.bot.enemy.Figures[figureID].Move(move.To.Position)
		//fmt.Println(move.bot.enemy.Figures)
	}
}
