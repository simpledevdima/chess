# chess
Chess game server

## Installation
get chess package
```
go install github.com/skvdmt/chess@latest
```

## Run chess server application
change directory to chess, build and run the application, then type in the browser:
Client application works on http://localhost:8080/chess
server on the addr ws://localhost:8081/ws/chess

## Future

server:
- chess bot development
- calculations about situations in which it is impossible to checkmate
- record of moves
- play record of moves

client:
- Display on the board captured pieces of the opponent
- Remember the move on the client and immediately do it when the move passes
- Check highlight
