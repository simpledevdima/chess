package chess_bot

func Start() {
	bot := NewBot() // making bot
	bot.setup()     // setting links and bot data
	bot.open()      // bot open connection to server
}
