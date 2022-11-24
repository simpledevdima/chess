package chess_bot

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/skvdmt/chess/game"
	"github.com/skvdmt/nrp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func NewBot() *Bot {
	return &Bot{
		ws:    "ws://localhost:8081/ws/chess",
		exit:  make(chan bool),
		send:  make(chan []byte),
		team:  &game.Team{},
		enemy: &game.Team{},
	}
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
	tpm       *game.TeamPossibleMoves
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
	//fmt.Println("team figures:", bot.team.Figures)
	//fmt.Println("enemy figures:", bot.enemy.Figures)
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

//func (bot *Bot) getRandomIndexMap(m interface{}) int {
//	keys := reflect.ValueOf(m).MapKeys()
//	return int(keys[rand.Intn(len(keys))].Int())
//}
//
//func (bot *Bot) getRandomMove() *move {
//	possibleMoves := *bot.team.GetPossibleMoves()
//	indexFigure := game.FigurerIndex(bot.getRandomIndexMap(possibleMoves))
//	indexMove := game.MoveIndex(bot.getRandomIndexMap(*possibleMoves[indexFigure]))
//	to := (*possibleMoves[indexFigure])[indexMove]
//	pos := bot.team.Figures[indexFigure].GetPosition()
//	return newMove(bot, pos, game.NewPosition(to.X, to.Y))
//}

func (bot *Bot) SetTeamPossibleMoves(tpm *game.TeamPossibleMoves) {
	bot.tpm = tpm
}

func (bot *Bot) move() {
	//randomMove := bot.getRandomMove()
	//randomMove.send()
	fmt.Println("MOVE")
	time.Sleep(time.Second / 10)

	bot.SetTeamPossibleMoves(bot.team.GetPossibleMoves())

	// set random rating
	for _, mvs := range *bot.tpm {
		for _, mv := range *mvs {
			mv.SetRating(rand.Float64() * 10)
		}
	}
	//bot.ShowPossibleMoves(bot.tpm)
	bot.team.ShowPossibleMoves(bot.tpm)

	// get move with max rating
	var rat float64
	var m *game.Move
	var fi game.FigurerIndex
	for index, mvs := range *bot.tpm {
		for _, mv := range *mvs {
			if mv.GetRating() > rat {
				rat = mv.GetRating()
				m = mv
				fi = index
			}
		}
	}
	bm := newMove(bot, bot.team.Figures[fi].GetPosition(), game.NewPosition(m.X, m.Y))
	bm.send()
}

// ShowPossibleMoves displays the possible moves of each piece of the team
//func (bot *Bot) ShowPossibleMoves(pm *game.TeamPossibleMoves) {
//	fmt.Printf("possible moves, team: %s\n", bot.team.Name.String())
//	for index, mvs := range *pm {
//		figure := bot.team.Figures[index]
//		x, y := figure.GetPosition().Get()
//		fmt.Printf("i=%2d n=%6s p=%dx%d to", index, figure.GetName(), x, y)
//		for _, mv := range *mvs {
//			fmt.Printf(" %dx%d(%.2f)", mv.X, mv.Y, mv.GetRating())
//		}
//		fmt.Printf("\n")
//	}
//	fmt.Println()
//}
