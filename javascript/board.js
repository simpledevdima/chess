class Board {
    constructor(chess) {
        // set link to chess
        this.chess = chess

        // making container
        this.container = document.createElement("canvas")

        // set board canvas context
        this.context = this.container.getContext("2d")

        // symbols cords of board
        this.symbols =  ["-", "a", "b", "c", "d", "e", "f", "g", "h"]

        // background colors
        this.backgroundColors = {
            back: "#ffffff",
            beige: "#ffce9e",
            brown: "#d18b47",
            beigeLastMove: "#ceff9e",
            brownLastMove: "#8bd147",
            beigeCheck: "#ff6666",
            brownCheck: "#cc3333",
        }

        // count of cells
        this.size = {
            x: 8,
            y: 8,
        }

        // one cell size
        this.cell = {
            width: 50,
            height: 50,
        }

        // teams
        this.teams = {
            white: {
                figures: [],
                eaten: [],
            },
            black: {
                figures: [],
                eaten: [],
            },
        }

        // pawn double move
        this.pawnDoubleMove = {}

        // figure taken from the board
        this.taken = {}

        // get image class
        this.images = new Images()

        this.setBoardContainer()

        // last move data
        this.lastMove = {}

        // move data
        this.move = {}

        // request animation frame ID
        this.requestAnimationFrameID = false
    }

    setBoardContainer() {
        // width and height of board
        this.container.width = 500
        this.container.height = 500

        // set class name
        this.container.classList.add("board")
    }

    // get show team name
    getShowTeamName() {
        if (this.chess.game.teamName === "spectators") {
            return "white"
        }
        return this.chess.game.teamName
    }

    // get show enemy name
    getShowEnemyName() {
        if (this.chess.game.enemyName === "") {
            return "black"
        }
        return this.chess.game.enemyName
    }

    // get notice cords
    calculateBoardCords(x, y) {
        switch (this.chess.game.teamName) {
            case "spectators":
            case "white":
                return [x, (this.size.y - y + 1)]
            case "black":
                return [(this.size.x - x + 1), y]
        }
    }

    // calculate board backlight color
    calculateBoardBacklightColor(x, y) {
        if ((x + y) / 2 - Math.floor((x + y) / 2) === 0) {
            return this.backgroundColors.beigeLastMove
        } else {
            return this.backgroundColors.brownLastMove
        }
    }

    // draw board last move backlight
    drawBoardLastMoveBacklight() {
        if (typeof this.lastMove.from !== "undefined") {
            const [fromX, fromY] = this.calculateBoardCords(this.lastMove.from.position.x, this.lastMove.from.position.y)
            this.context.fillStyle = this.calculateBoardBacklightColor(fromX, fromY)
            this.context.fillRect(fromX * this.cell.width, fromY * this.cell.height, this.cell.width, this.cell.height)

            const [toX, toY] = this.calculateBoardCords(this.lastMove.to.position.x, this.lastMove.to.position.y)
            this.context.fillStyle = this.calculateBoardBacklightColor(toX, toY)
            this.context.fillRect(toX * this.cell.width, toY * this.cell.height, this.cell.width, this.cell.height)
        }
    }

    // draw board on canvas context
    drawBoard() {
        this.context.fillStyle = this.backgroundColors.beige
        this.context.fillRect(50,50, this.cell.width * this.size.x, this.cell.height * this.size.y)
        this.context.fillStyle = this.backgroundColors.brown
        for (let x=1; x<=this.size.x; x++) {
            for (let y=1; y<=this.size.y; y++) {
                if (x / 2 - Math.trunc(x / 2) === 0 ^ y / 2 - Math.trunc(y / 2) === 0) {
                    this.context.fillRect(x * this.cell.width, y * this.cell.height, this.cell.width, this.cell.height)
                }
            }
        }
        this.drawBoardLastMoveBacklight()
    }

    // get horizontal symbol
    getNoticeHorizontalSymbol(x) {
        switch (this.chess.game.teamName) {
            case "spectators":
            case "white":
                return this.symbols[x].toString()
            case "black":
                return this.symbols[this.size.x - x + 1].toString()
        }
    }

    // get vertical notice
    getNoticeVerticalSymbol(y) {
        switch (this.chess.game.teamName) {
            case "spectators":
            case "white":
                return (this.size.y - y + 1).toString()
            case "black":
                return y.toString()
        }
    }

    // draw notices on canvas context
    drawNotices() {
        // background
        this.context.fillStyle = this.backgroundColors.back
        this.context.fillRect(0,0, this.container.width, this.container.height)

        this.context.fillStyle = "black"
        this.context.font = "18px Microsoft YaHei"
        for (let x=1; x<=this.size.x; x++) {
            this.context.fillText(this.getNoticeHorizontalSymbol(x), x * this.cell.width + 19, 32)
            this.context.fillText(this.getNoticeHorizontalSymbol(x), x * this.cell.width + 19, (this.size.y + 1) * this.cell.height + 32)
        }
        for (let y=1; y<=this.size.y; y++) {
            this.context.fillText(this.getNoticeVerticalSymbol(y), 19, y * this.cell.height + 32)
            this.context.fillText(this.getNoticeVerticalSymbol(y), (this.size.x + 1) * this.cell.width + 19, y * this.cell.height + 32)
        }
    }

    // calculate draw figure position
    calculateDrawFigurePosition(figure) {
        switch (this.chess.game.teamName) {
            case "spectators":
            case "white":
                return [
                    figure.position.x * this.cell.width + ((this.cell.width - this.images[this.getShowTeamName()][figure.name].width) / 2),
                    (this.size.y-figure.position.y+1) * this.cell.height + ((this.cell.height - this.images[this.getShowTeamName()][figure.name].height) / 2)
                ]
            case "black":
                return [
                    (this.size.x-figure.position.x+1) * this.cell.width + ((this.cell.width - this.images[this.getShowTeamName()][figure.name].width) / 2),
                    figure.position.y * this.cell.height + ((this.cell.height - this.images[this.getShowTeamName()][figure.name].height) / 2)
                ]
        }
    }

    // draw figures on board
    drawFigures() {
        for (const figure of Object.values(this.teams[this.getShowTeamName()].figures)) {
            const [x, y] = this.calculateDrawFigurePosition(figure)
            this.context.drawImage(this.images[this.getShowTeamName()][figure.name], x, y)
        }
        for (const figure of Object.values(this.teams[this.getShowEnemyName()].figures)) {
            const [x, y] = this.calculateDrawFigurePosition(figure)
            this.context.drawImage(this.images[this.getShowEnemyName()][figure.name], x, y)
        }
    }

    // get click on board cords
    getClickCords(event) {
        return [
            event.clientX - this.container.getBoundingClientRect().left,
            event.clientY - this.container.getBoundingClientRect().top
        ]
    }

    // get field cords
    getFieldCords(clickX, clickY) {
        for (let x = 1; x <= this.size.x; x++) {
            for (let y = 1; y <= this.size.y; y++) {
                if (clickX >= x * this.cell.width && clickX < (x + 1) * this.cell.width && clickY >= y * this.cell.height && clickY < (y + 1) * this.cell.height) {
                    return [x, y]
                }
            }
        }
        return [0, 0]
    }

    // calculation of coordinates for taking a figure in hand
    calculateForTakingFigure(x, y) {
        switch (this.chess.game.teamName) {
            case "white":
                return [x, this.size.y-y+1]
            case "black":
                return [this.size.x-x+1, y]
        }
    }

    // take the figure in hand
    takeTheFigureInHand(x, y) {
        for (const id of Object.keys(this.teams[this.chess.game.teamName].figures)) {
            const figure = this.teams[this.chess.game.teamName].figures[id]
            if (figure.position.x === x && figure.position.y === y) {
                this.chess.board.container.classList.add("cursor-"+this.chess.game.teamName+"-"+figure.name)
                this.taken.id = id
                this.taken.figure = figure
                delete this.teams[this.chess.game.teamName].figures[id]
            }
        }
    }

    // stay back taking figure
    stayBackTakingFigure() {
        this.teams[this.chess.game.teamName].figures[this.taken.id] = this.taken.figure
        this.chess.board.container.classList.remove("cursor-"+this.chess.game.teamName+"-"+this.taken.figure.name)
        this.taken = {}
    }

    // move taking figure
    moveTakingFigure(x, y) {
        this.chess.connection.sendRequest({
            post: "move",
            body: {
                from: {
                    position: {"x": this.taken.figure.position.x, "y": this.taken.figure.position.y}
                },
                to: {
                    position: {"x": x, "y": y}
                }
            }
        })
    }

    // use the figure from hand
    useTheFigureFromHand(x, y) {
        if (this.taken.figure.position.x === x && this.taken.figure.position.y === y) {
            this.stayBackTakingFigure()
        } else {
            this.moveTakingFigure(x, y)
        }
    }

    // click on board
    onClick() {
        const board = this
        this.onClickEvent = function (event) {
            const [clickX, clickY] = board.getClickCords(event)
            const [x, y] = board.getFieldCords(clickX, clickY)
            const [takeX, takeY] = board.calculateForTakingFigure(x, y)
            if (board.taken.id === undefined) {
                board.takeTheFigureInHand(takeX, takeY)
            } else {
                board.useTheFigureFromHand(takeX, takeY)
            }
        }
        this.chess.board.container.addEventListener("click", board.onClickEvent)
    }

    // transform pawn to queen
    transformPawnToQueen(figure) {
        if (figure.name === "pawn" && (figure.position.y === 1 || figure.position.y === 8)) {
            figure.name = "queen"
        }
    }

    getAnyFigureByCords(x, y) {
        return this.getFigureByCords(x, y, [this.getShowTeamName(), this.getShowEnemyName()])
    }

    getAnyFigureIDByCords(x, y) {
        return this.getFigureByCords(x, y, [this.getShowTeamName(), this.getShowEnemyName()], true)
    }

    getTeamFigureByCords(x, y) {
        return this.getFigureByCords(x,y,[this.getShowTeamName(), this.getShowEnemyName()], false, true )
    }

    // get figure by cords
    getFigureByCords(x, y, teams, returnID, returnTeam) {
        for (const teamName of teams) {
            for (const id of Object.keys(this.teams[teamName].figures)) {
                const figure = this.teams[teamName].figures[id]
                if (figure.position.x === x && figure.position.y === y) {
                    if (returnID) {
                        return id
                    } else if (returnTeam) {
                        return teamName
                    }
                    return figure
                }
            }
        }
        return false
    }

    // set move
    setMove(move) {
        this.move = move
    }

    // move figure from hand
    moveFigureFromHand() {
        this.taken.figure.position = this.move.to.position
        this.chess.board.container.classList.remove("cursor-"+this.chess.game.teamName+"-"+this.taken.figure.name)
        this.transformPawnToQueen(this.taken.figure)
        this.teams[this.chess.game.teamName].figures[this.taken.id] = this.taken.figure
        this.taken = {}
    }

    // move figure from board
    moveFigureFromBoard() {
        const figure = this.getAnyFigureByCords(this.move.from.position.x, this.move.from.position.y)
        if (figure) {
            figure.position = this.move.to.position
            this.transformPawnToQueen(figure)
        }
    }

    // eating enemy figure
    eatingFigure() {
        const eatenFigureTeam = this.getTeamFigureByCords(this.move.to.position.x, this.move.to.position.y)
        const eatenFigureID = this.getAnyFigureIDByCords(this.move.to.position.x, this.move.to.position.y)
        if (eatenFigureID) {
            this.teams[eatenFigureTeam].eaten[eatenFigureID] = this.teams[eatenFigureTeam].figures[eatenFigureID]
            delete this.teams[eatenFigureTeam].figures[eatenFigureID]
        }
    }

    // clear pawn double move
    clearPawnDoubleMove() {
        this.pawnDoubleMove = {}
    }

    // the pawn makes a double move
    pawnMakesDoubleMove() {
        this.clearPawnDoubleMove()
        const figure = this.getAnyFigureByCords(this.move.to.position.x, this.move.to.position.y)
        if (figure && figure.name === "pawn") {
            if (this.move.to.position.y === (this.move.from.position.y + 2 || this.move.from.position.y - 2)) {
                // add pawn double move
                this.pawnDoubleMove = {
                    id: this.getAnyFigureIDByCords(this.move.to.position.x, this.move.to.position.y),
                    teamName: this.getTeamFigureByCords(this.move.to.position.x, this.move.to.position.y),
                    x: this.move.to.position.x,
                    y: (this.move.to.position.y + this.move.from.position.y) / 2,
                }
            }
        }
        return false
    }

    // pawn take on the pass
    pawnTakeOnThePass() {
        const figure = this.getAnyFigureByCords(this.move.to.position.x, this.move.to.position.y)
        if (figure && figure.name === "pawn") {
            if (typeof this.pawnDoubleMove.id !== "undefined") {
                if (this.move.to.position.x === this.pawnDoubleMove.x && this.move.to.position.y === this.pawnDoubleMove.y) {
                    this.eatingFigure(
                        this.teams[this.pawnDoubleMove.teamName].figures[this.pawnDoubleMove.id].position.x,
                        this.teams[this.pawnDoubleMove.teamName].figures[this.pawnDoubleMove.id].position.y
                    )
                }
            }
        }
    }

    // last move backlight
    lastMoveBacklight() {
        this.lastMove = this.move
    }

    // exec move
    execMove() {
        this.chess.cause.hideCause()
        this.lastMoveBacklight()
        this.eatingFigure()
        if (typeof this.taken.id !== "undefined" && this.taken.figure.position.x === this.move.from.position.x && this.taken.figure.position.y === this.move.from.position.y) {
            this.moveFigureFromHand()
        } else {
            this.moveFigureFromBoard()
        }
        this.pawnTakeOnThePass()
        this.pawnMakesDoubleMove()
    }


    // set board figures
    setBoardFigures(board) {
        this.teams["white"].figures = board.white.figures
        this.teams["black"].figures = board.black.figures
        //console.log(this.teams)
        this.drawNotices()
        if (this.chess.game.youArePlayer) {
            this.onClick()
        }
        if (!this.requestAnimationFrameID) {
            this.animation()
        }
    }

    // refresh animation board
    animation() {
        this.requestAnimationFrameID = window.requestAnimationFrame(this.animation.bind(this))
        this.drawBoard()
        if (this.images.allImagesLoaded) {
            this.drawFigures()
        }
    }

}
