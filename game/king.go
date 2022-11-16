package game

func NewKing(x, y int, t *Team) *King {
	f := &King{}
	f.SetName("king")
	f.SetPosition(x, y)
	f.SetTeam(t)
	return f
}

// King is data type of chess figure
type King struct {
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (k *King) DetectionOfPossibleMove() []*Position {
	var possibleMoves []*Position
	for _, position := range k.detectionOfBrokenFields() {
		if !k.team.FigureExist(position.X, position.Y) && !k.kingOnTheBeatenFieldAfterMove(position.X, position.Y) {
			possibleMoves = append(possibleMoves, position)
		}
	}
	return possibleMoves
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (k *King) detectionOfBrokenFields() []*Position {
	var data []*Position

	if k.coordsOnBoard(k.X, k.Y+1) {
		data = append(data, NewPosition(k.X, k.Y+1))
	}
	if k.coordsOnBoard(k.X+1, k.Y+1) {
		data = append(data, NewPosition(k.X+1, k.Y+1))
	}
	if k.coordsOnBoard(k.X+1, k.Y) {
		data = append(data, NewPosition(k.X+1, k.Y))
	}
	if k.coordsOnBoard(k.X+1, k.Y-1) {
		data = append(data, NewPosition(k.X+1, k.Y-1))
	}
	if k.coordsOnBoard(k.X, k.Y-1) {
		data = append(data, NewPosition(k.X, k.Y-1))
	}
	if k.coordsOnBoard(k.X-1, k.Y-1) {
		data = append(data, NewPosition(k.X-1, k.Y-1))
	}
	if k.coordsOnBoard(k.X-1, k.Y) {
		data = append(data, NewPosition(k.X-1, k.Y))
	}
	if k.coordsOnBoard(k.X-1, k.Y+1) {
		data = append(data, NewPosition(k.X-1, k.Y+1))
	}

	return data
}

// Validation return true if this move are valid or return false
func (k *King) Validation(x int, y int) (bool, string) {
	if !k.coordsOnBoard(x, y) {
		return false, "attempt to go out the board"
	}
	if k.X == x && k.Y == y {
		return false, "can't walk around"
	}
	if k.team.FigureExist(x, y) {
		return false, "this place is occupied by your figure"
	}
	if k.kingOnTheBeatenFieldAfterMove(x, y) {
		return false, "your king stands on a beaten field"
	}
	// castling
	if !k.alreadyMove {
		if x == 3 {
			if !k.team.CheckingCheck() &&
				!k.team.FigureExist(k.X-1, k.Y) && !k.team.enemy.FigureExist(k.X-1, k.Y) &&
				!k.team.FigureExist(k.X-2, k.Y) && !k.team.enemy.FigureExist(k.X-2, k.Y) &&
				!k.team.FigureExist(k.X-3, k.Y) && !k.team.enemy.FigureExist(k.X-3, k.Y) &&
				k.team.FigureExist(k.X-4, k.Y) {
				if !k.team.GetFigureByCoords(k.X-4, k.Y).IsAlreadyMove() {
					return true, ""
				}
			}
		} else if x == 7 {
			if !k.team.CheckingCheck() &&
				!k.team.FigureExist(k.X+1, k.Y) && !k.team.enemy.FigureExist(k.X+1, k.Y) &&
				!k.team.FigureExist(k.X+2, k.Y) && !k.team.enemy.FigureExist(k.X+2, k.Y) &&
				k.team.FigureExist(k.X+3, k.Y) {
				if !k.team.GetFigureByCoords(k.X+3, k.Y).IsAlreadyMove() {
					return true, ""
				}
			}
		}
	}
	// detect Position for move and check it for input data move coords
	for _, position := range k.detectionOfBrokenFields() {
		if position.X == x && position.Y == y {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (k *King) Move(x int, y int) {
	k.team.pawnDoubleMove.clearPawnDoubleMove()
	k.MoveFigure(x, y)
}
