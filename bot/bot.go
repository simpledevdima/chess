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
	"reflect"
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
			m := NewMove(bot, nil, nil)
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

func (bot *Bot) getRandomMove() *move {
	possibleMoves := bot.team.GetPossibleMoves()
	figuresKeys := reflect.ValueOf(possibleMoves).MapKeys()
	index := game.FigureIndex(int(figuresKeys[rand.Intn(len(figuresKeys))].Int()))
	to := possibleMoves[index][rand.Intn(len(possibleMoves[index]))]
	fx, fy := bot.team.Figures[index].GetPosition()
	return NewMove(bot, game.NewPosition(fx, fy), game.NewPosition(to.X, to.Y))
}

func (bot *Bot) move() {
	//bot.ShowMoves(bot.team)
	//bot.ShowMoves(bot.enemy)
	time.Sleep(time.Second / 10)
	randomMove := bot.getRandomMove()
	randomMove.send()
}

//func (bot *Bot) ShowMoves(t *game.Team) {
//	fmt.Println(t.Name.String())
//	for index, figure := range t.Figures {
//		x, y := figure.GetPosition()
//		moves := figure.DetectionOfPossibleMove()
//		fmt.Println(index, figure.GetName(), x, y)
//		for _, m := range moves {
//			x, y := m.Get()
//			fmt.Print(x, y, " | ")
//		}
//		fmt.Println()
//	}
//}

//
//func (bot *Bot) sendTestRequest() {
//	// test send data
//	type request struct{
//		Id int
//		Post string
//		Body interface{}
//	}
//	var req request
//	req.Id = 328748324
//	req.Post = "setup"
//	req.Body = struct{}{}
//	dataJSON, err := json.Marshal(req)
//	if err != nil {
//		log.Println(err)
//	}
//	fmt.Println(string(dataJSON))
//	bot.send <- dataJSON
//}

//fmt.Println(possibleMoves)
//fmt.Println("figureID", index)
//fmt.Println("from", fx, fy)
//fmt.Println("to", to.X, to.Y)

//fmt.Println(bot.data.isAllDataLoaded())
//fmt.Println(bot.data.isPlay())
//fmt.Println(!bot.data.isOver())
//fmt.Println(bot.data.yourTurn())
//if bot.data.isPlay() && !bot.data.isOver() && bot.data.yourTurn() {
//	fmt.Println("game play")
//	//bot.data.move()
//} else {
//	fmt.Println("game stop")
//}
