class Scoreboard {
    constructor(chess) {
        this.chess = chess

        // making scoreboard
        this.container = document.createElement("div")
        this.container.classList.add("scoreboard")

        // making elements
        this.team = document.createElement("div")
        this.teamNotice = document.createElement("span")
        this.teamName = document.createElement("span")
        this.condition = document.createElement("div")
        this.notice = document.createElement("div")
        this.turn = document.createElement("div")
        this.timeLeft = document.createElement("div")
        this.reserveTimeLeft = document.createElement("div")
        this.reserveTimeLeftHeader = document.createElement("div")
        this.reserveTimeLeftWhite = document.createElement("div")
        this.reserveTimeLeftBlack = document.createElement("div")

        // making scoreboard
        this.container.appendChild(this.team)
        this.team.appendChild(this.teamNotice)
        this.team.appendChild(this.teamName)
        this.container.appendChild(this.condition)
        this.container.appendChild(this.notice)
        this.container.appendChild(this.turn)
        this.container.appendChild(this.timeLeft)
        this.container.appendChild(this.reserveTimeLeft)
        this.reserveTimeLeft.appendChild(this.reserveTimeLeftHeader)
        this.reserveTimeLeft.appendChild(this.reserveTimeLeftWhite)
        this.reserveTimeLeft.appendChild(this.reserveTimeLeftBlack)

        this.setupScoreboardElements()
    }

    setupScoreboardElements() {
        // set classes
        this.team.classList.add("team")
        this.teamNotice.classList.add("notice")
        this.teamName.classList.add("name")
        this.condition.classList.add("condition")
        this.notice.classList.add("notice")
        this.turn.classList.add("turn")
        this.timeLeft.classList.add("time-left")
        this.reserveTimeLeft.classList.add("reserve-time-left")
        this.reserveTimeLeftHeader.classList.add("header")
        this.reserveTimeLeftWhite.classList.add("white")
        this.reserveTimeLeftBlack.classList.add("black")

        // set content
        this.teamNotice.innerHTML = "your team:"
        this.timeLeft.innerHTML = "&nbsp;"
        this.reserveTimeLeftHeader.innerHTML = "reserve time:"
    }

    showTeamName() {
        this.teamName.innerHTML = this.chess.game.teamName
    }

    // why game are stopped set notice
    whyGameStoppedNotice() {
        this.notice.innerHTML = this.chess.game.playCause
    }

    // show play
    showPlay() {
        this.turn.classList.remove("over")
        this.condition.classList.remove("over")
        if (this.chess.game.play) {
            this.condition.innerHTML = "Game played"
            this.notice.innerHTML = "&nbsp;"
            this.changeTurnBacklight()
        } else {
            this.condition.innerHTML = "Game stopped"
            this.whyGameStoppedNotice()
        }
    }

    // set win notice
    setWinNotice() {
        this.notice.innerHTML = this.chess.game.overCause
    }

    // show over
    showOver() {
        this.turn.classList.add("over")
        this.condition.classList.add("over")
        this.condition.innerHTML = "GAME OVER"
        this.turn.classList.remove("your-move")
        this.setWinNotice()
    }

    // show turn
    showTurn(turn) {
        this.turn.innerHTML = turn + " turn"
    }

    // change turn backlight
    changeTurnBacklight() {
        this.turn.classList.remove("your-move")
        if (this.chess.game.play && ((this.chess.game.teamName === "white" && this.chess.game.turn === "white") || (this.chess.game.teamName === "black" && this.chess.game.turn === "black"))) {
            this.turn.classList.add("your-move")
        }
    }

    // get seconds int convert and return time in MM:SS format
    secondsToMMSS(seconds) {
        let mm = Math.floor(seconds / 60);
        let ss = seconds - (mm * 60);
        if (ss < 10) {ss = "0"+ss;}
        return mm + ':' + ss;
    }

    // time left change backlight
    timeLeftChangeBacklight() {
        if (this.chess.game.time.left <= 5) {
            if (this.chess.game.time.teamName === this.chess.game.teamName) {
                this.timeLeft.classList.add("warning")
            }
        } else {
            this.timeLeft.classList.remove("warning")
        }
    }

    // show step time left
    showStepTimeLeft() {
        this.timeLeft.innerHTML = this.secondsToMMSS(this.chess.game.time.left) + " left"
    }

    reserveTimeLeftChangeTimers() {
        switch (this.chess.game.time.teamName) {
            case "white":
                this.reserveTimeLeftWhite.innerHTML = this.secondsToMMSS(this.chess.game.time.left) + " white"
                break
            case "black":
                this.reserveTimeLeftBlack.innerHTML = this.secondsToMMSS(this.chess.game.time.left) + " black"
                break
        }
    }

    showReserveTimeLeft() {
        this.timeLeft.innerHTML = this.secondsToMMSS(this.chess.game.time.left) + " left"
    }
}
