// Game class for saving and management game data and event handling
class Game {
    constructor(chess) {
        this.chess = chess

        // your team name
        this.teamName = ""

        // your enemy team name
        this.enemyName = ""

        // game are play
        this.play = false
        this.playCause = ""

        // game are over
        this.over = false
        this.overCause = ""

        // turn
        this.turn = ""

        // time left
        this.time = {
            left: 0,
            teamName: "",
        }
    }

    // set time
    setTime(left, teamName) {
        this.time.left = left
        this.time.teamName = teamName
    }

    // set turn
    setTurn(turn) {
        this.turn = turn
    }

    // set game over
    setOver(data) {
        this.over = data.over
        if (typeof data.cause) {
            this.overCause = data.cause
        }
    }

    // set game play
    setPlay(data) {
        this.play = data.play
        if (typeof data.cause !== "undefined") {
            this.playCause = data.cause
        } else {
            this.playCause = ""
        }
    }

    // set team name
    setTeamName(name) {
        this.teamName = name
        if (this.teamName === "spectators") {
            this.chess.actions.hideActions()
        }
    }

    // set enemy name opposite team name
    setEnemyName() {
        switch (this.teamName) {
            case "white":
                this.enemyName = "black"
                break
            case "black":
                this.enemyName = "white"
                break
        }
    }

    get youArePlayer() {
        return this.teamName === "white" || this.teamName === "black";
    }

    setStart() {
        this.chess.board.lastMove = {}
    }

    execOver() {
        this.chess.actions.offerDrawDecisions.classList.add("hide")
        this.chess.board.container.removeEventListener("click" , this.chess.board.onClickEvent)
    }
}