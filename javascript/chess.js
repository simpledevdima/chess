document.addEventListener('DOMContentLoaded', function() {
    // get chess object
    const chess = document.querySelector(".chess")

    // making containers
    chess.containerLeft = document.createElement("div")
    chess.containerRight = document.createElement("div")
    chess.containerRight.classList.add("right")

    // set classes
    chess.config = new Config()
    chess.board = new Board(chess)
    chess.cause = new Cause(chess)
    chess.scoreboard = new Scoreboard(chess)
    chess.actions = new Actions(chess)
    chess.game = new Game(chess)
    chess.connection = new Connection(chess)

    // making chess
    chess.appendChild(chess.containerLeft)
    chess.containerLeft.appendChild(chess.board.container)
    chess.containerLeft.appendChild(chess.cause.container)
    chess.appendChild(chess.containerRight)
    chess.containerRight.appendChild(chess.scoreboard.container)
    chess.containerRight.appendChild(chess.actions.container)
})
