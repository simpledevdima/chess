package game

import (
	"log"
)

// King is data type of chess figure
type King struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (king *King) DetectionOfPossibleMove() []Position {
	var possibleMoves []Position
	for _, position := range king.detectionOfBrokenFields() {
		if !king.team.FigureExist(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (king *King) detectionOfBrokenFields() []Position {
	var data []Position

	if king.coordsOnBoard(king.X, king.Y+1) {
		data = append(data, Position{X: king.X, Y: king.Y + 1})
	}
	if king.coordsOnBoard(king.X+1, king.Y+1) {
		data = append(data, Position{X: king.X + 1, Y: king.Y + 1})
	}
	if king.coordsOnBoard(king.X+1, king.Y) {
		data = append(data, Position{X: king.X + 1, Y: king.Y})
	}
	if king.coordsOnBoard(king.X+1, king.Y-1) {
		data = append(data, Position{X: king.X + 1, Y: king.Y - 1})
	}
	if king.coordsOnBoard(king.X, king.Y-1) {
		data = append(data, Position{X: king.X, Y: king.Y - 1})
	}
	if king.coordsOnBoard(king.X-1, king.Y-1) {
		data = append(data, Position{X: king.X - 1, Y: king.Y - 1})
	}
	if king.coordsOnBoard(king.X-1, king.Y) {
		data = append(data, Position{X: king.X - 1, Y: king.Y})
	}
	if king.coordsOnBoard(king.X-1, king.Y+1) {
		data = append(data, Position{X: king.X - 1, Y: king.Y + 1})
	}

	return data
}

// Validation return true if this move are valid or return false
func (king *King) Validation(x int, y int) (bool, string) {
	if !king.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if king.X == x && king.Y == y {
		return false, "can't walk around"
	}
	if king.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if king.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// castling
	if !king.alreadyMove {
		if x == 3 {
			if !king.team.CheckingCheck() &&
				!king.team.FigureExist(king.X-1, king.Y) && !king.enemy.FigureExist(king.X-1, king.Y) &&
				!king.team.FigureExist(king.X-2, king.Y) && !king.enemy.FigureExist(king.X-2, king.Y) &&
				!king.team.FigureExist(king.X-3, king.Y) && !king.enemy.FigureExist(king.X-3, king.Y) &&
				king.team.FigureExist(king.X-4, king.Y) {
				figureID, err := king.team.GetFigureID(king.X-4, king.Y)
				if err != nil {
					log.Println(err)
				}
				if !king.team.Figures[figureID].IsAlreadyMove() {
					return true, ""
				}
			}
		} else if x == 7 {
			if !king.team.CheckingCheck() &&
				!king.team.FigureExist(king.X+1, king.Y) && !king.enemy.FigureExist(king.X+1, king.Y) &&
				!king.team.FigureExist(king.X+2, king.Y) && !king.enemy.FigureExist(king.X+2, king.Y) &&
				king.team.FigureExist(king.X+3, king.Y) {
				figureID, err := king.team.GetFigureID(king.X+3, king.Y)
				if err != nil {
					log.Println(err)
				}
				if !king.team.Figures[figureID].IsAlreadyMove() {
					return true, ""
				}
			}
		}
	}
	// detect Position for move and check it for input data move coords
	for _, position := range king.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (king *King) Move(x int, y int) {
	king.team.pawnDoubleMove.clearPawnDoubleMove()
	king.MoveFigure(x, y)
}
