package chess_bot

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/simpledevdima/chess/game"
	"github.com/simpledevdima/nrp"
	"log"
	"net/http"
	"time"
)

func NewBot() *Bot {
	b := &Bot{
		ws:     "ws://localhost:8081/ws/chess",
		exit:   make(chan bool),
		send:   make(chan []byte),
		team:   &game.Team{},
		enemy:  &game.Team{},
		rating: newRating(),
	}
	b.rating.setBot(b)
	return b
}

type Bot struct {
	ws        string
	wsHeaders http.Header
	exit      chan bool
	send      chan []byte
	conn      *websocket.Conn
	status    status
	team      *game.Team
	enemy     *game.Team
	rating    *rating
	//board     *board
}

func (bot *Bot) setLinks() {
	bot.team.SetEnemy(bot.enemy)
	bot.enemy.SetEnemy(bot.team)
	bot.status.setBot(bot)
}

// setup data bot
func (bot *Bot) setup() {
	bot.setLinks()
	bot.setWSHeaders()
	bot.team.MakeFigures()
	bot.enemy.MakeFigures()
}

// setBoard set white and black team figures data
func (bot *Bot) setBoard(request *nrp.Simple) {
	var body struct {
		White struct {
			Figures interface{} `json:"figures"`
		} `json:"white"`
		Black struct {
			Figures interface{} `json:"figures"`
		} `json:"black"`
	}
	request.BodyToVariable(&body)
	whiteFigures, err := json.Marshal(body.White.Figures)
	if err != nil {
		log.Println(err)
	}
	bot.setFigures("white", whiteFigures)
	blackFigures, err := json.Marshal(body.Black.Figures)
	if err != nil {
		log.Println(err)
	}
	bot.setFigures("black", blackFigures)
}

func (bot *Bot) setFigures(team string, figures []byte) {
	switch {
	case team == bot.status.teamName:
		bot.team.ImportFigures(figures)
	case team != bot.status.teamName:
		bot.enemy.ImportFigures(figures)
	}
}

// setWSHeaders setup http headers for ws connection
func (bot *Bot) setWSHeaders() {
	bot.wsHeaders = make(http.Header)
	bot.wsHeaders.Set("Origin", "http://localhost:8080")
}

func (bot *Bot) open() {
	fmt.Printf("connecting to %s\n", bot.ws)
	var err error
	bot.conn, _, err = websocket.DefaultDialer.Dial(bot.ws, bot.wsHeaders)
	if err != nil {
		log.Println(err)
	}
	defer bot.close()
	go bot.read()
	go bot.write()
	bot.wait()
}

func (bot *Bot) close() {
	fmt.Println("connection close")
	err := bot.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (bot *Bot) exitApp() {
	bot.exit <- true
}

func (bot *Bot) read() {
	defer close(bot.send)
	for {
		_, dataJSON, err := bot.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			bot.exitApp()
			break
		}
		//fmt.Printf("recv: %s\n", dataJSON)

		request := nrp.Simple{}
		request.Parse(dataJSON)
		//fmt.Println(string(request.Export()))
		switch request.Post {
		case "your":
			bot.status.setTeamName(&request)
			bot.status.setEnemyName()
			// set team and enemy names
			bot.team.SetName(bot.status.teamName)
			bot.enemy.SetName(bot.status.enemyName)
		case "move":
			m := newMove(bot, nil, nil)
			request.BodyToVariable(&m)
			m.exec()
		case "board":
			bot.setBoard(&request)
		case "turn":
			bot.status.setTurn(&request)
			bot.status.checkForChanges()
		case "game_play":
			bot.status.setPlay(&request)
			bot.status.checkForChanges()
		case "game_over":
			bot.status.setOver(&request)
			bot.status.checkForChanges()
		default:
			//log.Println("unknown data")
		}
	}
}

func (bot *Bot) write() {
	for {
		select {
		case data := <-bot.send:
			err := bot.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				bot.exitApp()
			}
		}
	}
}

func (bot *Bot) wait() {
	for {
		select {
		case exit := <-bot.exit:
			if exit {
				return
			}
		}
	}
}

func (bot *Bot) move() {
	time.Sleep(time.Second / 10)

	bot.rating.setTeamPossibleMoves(bot.team.GetPossibleMoves())
	bot.rating.setTeamBrokenFields(bot.team.GetBrokenFields())
	bot.rating.setEnemyPossibleMoves(bot.enemy.GetPossibleMoves())
	bot.rating.setEnemyBrokenFields(bot.enemy.GetBrokenFields())

	bot.rating.setRandomRatingToPossibleMoves()

	//bot.board = newBoard()
	//bot.board.calculate()

	bot.rating.EatUnprotectedFigure()
	bot.rating.MoveToBrokenField()
	bot.rating.moveUnprotectedFigureFromBrokenFieldToProtectedOrSecureField()
	bot.rating.protectAttackedUnprotectedFigure()

	mv := bot.rating.getMoveWithMaxRating()
	mv.send()
}
