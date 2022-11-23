package game

import (
	"errors"
	"fmt"
	"log"
)

// Figure information about each figure
type Figure struct {
	figurer     Figurer
	Name        string `json:"name"`
	*Position   `json:"position"`
	alreadyMove bool
	team        *Team
}

// ErrorDetail desc
func (f *Figure) ErrorDetail(pos *Position) error {
	switch {
	case !f.positionOnBoard(pos):
		return errors.New("attempt to go out the board")
	case *f.GetPosition() == *pos:
		return errors.New("can't walk around")
	case !f.figurer.CanWalkLikeThat(pos):
		return errors.New(fmt.Sprintf("%s doesn't walk like that", f.GetName()))
	case f.team.Figures.ExistsByPosition(pos):
		return errors.New("this place is occupied by your figure")
	case f.kingOnTheBeatenFieldAfterMove(pos):
		return errors.New("your king stands on a beaten field")
	default:
		return errors.New("this figure cant make that move")
	}
}

// Validation return true if this move are valid or return false
func (f *Figure) Validation(pos *Position) (bool, error) {
	for _, position := range *f.figurer.GetPossibleMoves(false) {
		if *position == *pos {
			return true, nil
		}
	}
	return false, f.ErrorDetail(pos)
}

type Direction int

const (
	top Direction = iota
	topRight
	right
	rightBottom
	bottom
	bottomLeft
	left
	leftTop
)

func (f *Figure) GetPositionsByDirectionsAndMaxRemote(opened map[Direction]bool, maxRemote uint8) *Positions {
	poss := make(Positions)
	var pi PositionIndex
	for remote := uint8(1); remote <= maxRemote; remote++ {
		for dir := range opened {
			if opened[dir] {
				func() {
					pos := f.GetPositionByDirectionAndRemote(dir, remote)
					if !f.positionOnBoard(pos) {
						opened[dir] = false
						return
					}
					if f.team.Figures.ExistsByPosition(pos) ||
						f.team.enemy.Figures.ExistsByPosition(pos) {
						opened[dir] = false
					}
					pi = poss.Set(pi, pos)
				}()
			}
		}
	}
	return &poss
}

func (f *Figure) GetPositionByDirectionAndRemote(dir Direction, remote uint8) *Position {
	var pos *Position
	switch dir {
	case top:
		pos = NewPosition(f.X, f.Y+remote)
	case topRight:
		pos = NewPosition(f.X+remote, f.Y+remote)
	case right:
		pos = NewPosition(f.X+remote, f.Y)
	case rightBottom:
		pos = NewPosition(f.X+remote, f.Y-remote)
	case bottom:
		pos = NewPosition(f.X, f.Y-remote)
	case bottomLeft:
		pos = NewPosition(f.X-remote, f.Y-remote)
	case left:
		pos = NewPosition(f.X-remote, f.Y)
	case leftTop:
		pos = NewPosition(f.X-remote, f.Y+remote)
	}
	return pos
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (f *Figure) GetPossibleMoves(thereIs bool) *Positions {
	poss := make(Positions)
	var pi PositionIndex
	for _, position := range *f.figurer.GetBrokenFields() {
		if !f.team.Figures.ExistsByPosition(position) && !f.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
			if thereIs {
				return &poss
			}
		}
	}
	return &poss
}

// GetName get name from figure
func (f *Figure) GetName() string {
	return f.Name
}

// SetName set name to the figure
func (f *Figure) SetName(name string) {
	f.Name = name
}

// IsAlreadyMove returns true if the figure has moved in the current match otherwise returns false
func (f *Figure) IsAlreadyMove() bool {
	return f.alreadyMove
}

// setAlreadyMove setting a value about moving a figure in the current match
func (f *Figure) setAlreadyMove(flag bool) {
	f.alreadyMove = flag
}

// SimulationMove makes a move, executes a callback and returns to the starting situation
func (f *Figure) SimulationMove(pos *Position, callback func() bool) bool {
	// move
	curPos := f.GetPosition()
	f.SetPosition(pos)
	// undo move
	defer func() {
		f.SetPosition(curPos)
	}()
	// eating
	if f.team.enemy.Figures.ExistsByPosition(pos) {
		index, figure := f.team.enemy.Figures.GetIndexAndFigureByPosition(pos)
		f.team.enemy.Eaten.Set(index, figure)
		f.team.enemy.Figures.Remove(index)
		// undo eating
		defer func() {
			if f.team.enemy.Eaten.ExistsByPosition(pos) {
				f.team.enemy.Figures.Set(index, figure)
				f.team.enemy.Eaten.Remove(index)
			}
		}()
	}
	// take on the pass
	if f.team.enemy.pawnDoubleMove.isTakeOnThePass(pos) && f.GetName() == "pawn" {
		index, figure := f.team.enemy.Figures.GetIndexAndFigureByPosition(f.team.enemy.pawnDoubleMove.pawn.GetPosition())
		f.team.enemy.Eaten.Set(index, figure)
		f.team.enemy.Figures.Remove(index)
		// undo take on the pass
		defer func() {
			f.team.enemy.Figures.Set(index, figure)
			f.team.enemy.Eaten.Remove(index)
		}()
	}
	return callback()
}

// kingOnTheBeatenFieldAfterMove returns true if your king is on the beaten square after the move otherwise return false
func (f *Figure) kingOnTheBeatenFieldAfterMove(pos *Position) bool {
	return f.SimulationMove(pos, func() bool {
		return f.team.CheckingCheck()
	})
}

// positionOnBoard returns true if the coordinates are within the board, otherwise return false
func (f *Figure) positionOnBoard(pos *Position) bool {
	x, y := pos.Get()
	if x >= 1 && x <= 8 && y >= 1 && y <= 8 {
		return true
	}
	return false
}

// SetTeam set links to your team for current figure
func (f *Figure) SetTeam(team *Team) {
	f.team = team
}

// Move to new position and eat enemy figure if you need that
func (f *Figure) Move(pos *Position) {
	// pawn double move
	if f.GetName() == "pawn" {
		f.team.enemy.pawnDoubleMove.pawnTakeOnThePass(pos)
		f.team.pawnDoubleMove.pawnMakesDoubleMove(f, f.GetPosition(), pos)
	} else {
		f.team.pawnDoubleMove.clearPawnDoubleMove()
	}

	f.SetPosition(pos)
	f.setAlreadyMove(true)
	if f.team.enemy != nil && f.team.enemy.Figures.ExistsByPosition(pos) {
		// eat enemy figure
		err := f.team.enemy.Eating(pos)
		if err != nil {
			log.Println(err)
		}
	}

	f.transformPawnTOQueen(pos)

	// Debug
	//f.team.ShowBrokenFields()
	//f.team.enemy.ShowBrokenFields()
	//f.team.ShowPossibleMoves()
	//f.team.enemy.ShowPossibleMoves()
}

// transformPawnToQueen promote a pawn to a queen
func (f *Figure) transformPawnTOQueen(pos *Position) {
	if f.GetName() == "pawn" && (f.Y == 1 || f.Y == 8) {
		figureID := f.team.Figures.GetIndexByPosition(pos)
		// replace pawn to queen
		f.team.Figures.Set(figureID, NewQueen(pos, f.team))
		f.team.Figures.Get(figureID).setAlreadyMove(true)
	}
}

// SetPosition set Position to coords from argument
func (f *Figure) SetPosition(p *Position) {
	f.Position = p
}

// GetPosition return current Position of Figure
func (f *Figure) GetPosition() *Position {
	return f.Position
}
