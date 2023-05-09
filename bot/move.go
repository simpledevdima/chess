package chess_bot

import (
	"fmt"
	"github.com/simpledevdima/chess/game"
	"github.com/simpledevdima/nrp"
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
func (m *move) setBot(bot *Bot) {
	m.bot = bot
}

// send отправка данных о совершаемом ходе на сервер
func (m *move) send() {
	if m.bot.status.isYourTurn() {
		request := nrp.Simple{Post: "move", Body: m}
		fmt.Println("BOT SEND MOVE:", string(request.Export()))
		m.bot.send <- request.Export()
	}
}

// exec выполнение хода информаци о котором поступила с сервера
func (m *move) exec() {
	if m.bot.team.Figures.ExistsByPosition(m.From.Position) {
		// your team move
		figureID := m.bot.team.Figures.GetIndexByPosition(m.From.Position)
		m.bot.team.Figures[figureID].Move(m.To.Position)
	} else if m.bot.enemy.Figures.ExistsByPosition(m.From.Position) {
		// enemy team move
		figureID := m.bot.enemy.Figures.GetIndexByPosition(m.From.Position)
		m.bot.enemy.Figures[figureID].Move(m.To.Position)
	}
}
