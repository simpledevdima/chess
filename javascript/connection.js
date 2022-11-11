class Connection {
    constructor(chess) {
        this.chess = chess

        // requests
        this.requests = {}

        // create websocket connection
        this.socket = new WebSocket(this.chess.config.chessWSServerAddr);

        // start websocket events
        this.webSocketEvents()
    }

    // websocket events
    webSocketEvents() {
        const conn = this

        // websocket connection are close
        this.socket.onclose = function () {
            alert("Connection to server lost")
        }

        // came message
        this.socket.onmessage = function (answer) {
            const data = JSON.parse(answer.data);
            // console.log(data)
            if (typeof data.id == "undefined") {
                // event
                conn.onEvent(data)
            } else {
                // response to request
                conn.onResponse(data)
            }
        }
    }

    // send request to server and
    // save unique id and time of create to chess.requests
    // to process the server response
    sendRequest(request) {
        let type
        if (typeof request.post !== "undefined") {
            type = request.post
        } else {
            console.log("ERROR: Bad request")
            console.log(request)
            return false
        }
        do {
            request.id = Math.floor(Math.random() * 1000000)
        } while (typeof this.requests[request.id] !== "undefined")
        this.requests[request.id] = {dt_create: Math.floor(Date.now() / 1000), type: type}
        // console.log(request)
        this.socket.send(JSON.stringify(request))
        return true
    }

    onEvent(data) {
        if (typeof data.your_team_name !== "undefined") {
            this.chess.game.setTeamName(data.your_team_name)
            this.chess.game.setEnemyName()
            this.chess.scoreboard.showTeamName()
        }
        if (typeof data.board !== "undefined") {
            this.chess.board.setBoardFigures(data.board)
        }
        if (typeof data.play !== "undefined") {
            this.chess.game.setPlay(data)
            if (!this.chess.game.over) {
                this.chess.scoreboard.showPlay()
            }
        }
        if (typeof data.over !== "undefined") {
            this.chess.game.setOver(data)
            if (this.chess.game.over) {
                this.chess.game.execOver()
                this.chess.scoreboard.showOver()
                this.chess.actions.showOver()
            } else {
                this.chess.game.setStart()
                this.chess.scoreboard.showPlay()
                this.chess.actions.showPlay()
            }
        }
        if (typeof data.turn !== "undefined") {
            this.chess.game.setTurn(data.turn)
            this.chess.scoreboard.showTurn(data.turn)
            this.chess.scoreboard.changeTurnBacklight()
        }
        if (typeof data.step_time_left !== "undefined" && typeof data.team_name !== "undefined") {
            this.chess.game.setTime(data.step_time_left, data.team_name)
            this.chess.scoreboard.showStepTimeLeft()
            this.chess.scoreboard.timeLeftChangeBacklight()
        }
        if (typeof data.reserve_time_left !== "undefined" && typeof data.team_name !== "undefined") {
            this.chess.game.setTime(data.reserve_time_left, data.team_name)
            this.chess.scoreboard.showReserveTimeLeft()
            this.chess.scoreboard.reserveTimeLeftChangeTimers()
        }
        if (typeof data.move !== "undefined") {
            this.chess.board.setMove(data.move)
            this.chess.board.execMove()
        }
        if (typeof data.opponent_offer_a_draw !== "undefined") {
            this.chess.actions.showOfferDrawDecisions()
        }
        if (typeof data.time_left_for_confirm_draw !== "undefined") {
            this.chess.actions.changeOfferDrawTimer(data.time_left_for_confirm_draw)
        }
        if (typeof data.attempts_left_to_offer_a_draw !== "undefined") {
            this.chess.actions.setAttemptsLeftToOfferADraw(data.attempts_left_to_offer_a_draw)
            this.chess.actions.changeAttemptsLeftToOfferADrawButton()
        }
    }

    // move response processing
    moveResponseProcessing(data) {
        if (!data.body.valid) {
            // move not valid stay back taking figure
            this.chess.board.stayBackTakingFigure()
            this.chess.cause.showCause(data.body.cause)
        }
    }

    offerDrawResponse(data) {
        this.chess.actions.offerDrawWaitDecision = false
        this.chess.actions.offerDraw.disabled = false
        this.chess.actions.changeAttemptsLeftToOfferADrawButton()
        if (!data.body.valid) {
            this.chess.cause.showCause(data.body.cause)
        }
    }

    // server response to request
    onResponse(data) {
        switch (this.requests[data.id].type) {
            case "move":
                this.moveResponseProcessing(data)
                break
            case "offer_a_draw":
                this.offerDrawResponse(data)
                break
        }
    }

}
