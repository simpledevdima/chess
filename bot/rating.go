package chess_bot

import (
	"github.com/simpledevdima/chess/game"
	"math/rand"
)

// NewRating returns a link to a new rating
func newRating() *rating {
	r := &rating{}
	return r
}

// rating data type containing methods for calculating the rating of the team's possible moves
type rating struct {
	teamPossibleMoves  *game.TeamPossibleMoves
	teamBrokenFields   *game.TeamBrokenFields
	enemyPossibleMoves *game.TeamPossibleMoves
	enemyBrokenFields  *game.TeamBrokenFields
	bot                *Bot
}

func (r *rating) setBot(b *Bot) {
	r.bot = b
}

// setTeamPossibleMoves setting the value of the possible moves of the command from the argument
func (r *rating) setTeamPossibleMoves(tpms *game.TeamPossibleMoves) {
	r.teamPossibleMoves = tpms
}

// setTeamBrokenFields setting the value of broken command fields from the argument
func (r *rating) setTeamBrokenFields(tbfs *game.TeamBrokenFields) {
	r.teamBrokenFields = tbfs
}

// setEnemyPossibleMoves setting the value of the possible moves of the enemy team from the argument
func (r *rating) setEnemyPossibleMoves(epms *game.TeamPossibleMoves) {
	r.enemyPossibleMoves = epms
}

// setEnemyBrokenFields setting the value of the broken fields of the enemy team from the argument
func (r *rating) setEnemyBrokenFields(ebfs *game.TeamBrokenFields) {
	r.enemyBrokenFields = ebfs
}

// setRandomRatingToPossibleMoves sets a random rating from 0 to 1 for all possible moves
func (r *rating) setRandomRatingToPossibleMoves() {
	for _, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			mv.SetRating(rand.Float64())
		}
	}
	r.bot.team.ShowPossibleMoves(r.teamPossibleMoves)
}

// EatUnprotectedFigure метод изменяющий рейтинг хода на основе того, что фигура может съесть незащищенную фигуру противника
func (r *rating) EatUnprotectedFigure() {
	for _, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			if r.bot.enemy.Figures.ExistsByPosition(mv.Position) && !r.posIsAttacked(mv.Position) {
				ef := r.bot.enemy.Figures.GetByPosition(mv.Position)
				rat := r.getFigureValueByName(ef.GetName())
				mv.SetRating(mv.GetRating() + rat)
			}
		}
	}
}

// MoveToBrokenField changing the rating for moves to broken opponent's fields
func (r *rating) MoveToBrokenField() {
	r.teamPossibleMoves.IterMoves(func(_ game.MoveIndex, mv *game.Move) {
	})
	for fi, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			if r.posIsAttacked(mv.Position) {
				rat := r.getFigureValueByName(r.bot.team.Figures[fi].GetName())
				mv.SetRating(mv.GetRating() - rat)
			}
		}
	}
}

// remove an unprotected piece from the opponent's beaten square
func (r *rating) moveUnprotectedFigureFromBrokenFieldToProtectedOrSecureField() {
	for fi, mvs := range *r.teamPossibleMoves {
		if !r.posIsProtected(r.bot.team.Figures[fi].GetPosition()) && r.posIsAttacked(r.bot.team.Figures[fi].GetPosition()) {
			for _, mv := range *mvs {
				if (r.posIsProtected(mv.Position) && r.posIsAttacked(mv.Position)) || !r.posIsAttacked(mv.Position) {
					rat := r.getFigureValueByName(r.bot.team.Figures[fi].GetName())
					mv.SetRating(mv.GetRating() + rat)
				}
			}
		}
	}
}

// защита незащищенных атакуемых соперником фигур
func (r *rating) protectAttackedUnprotectedFigure() {
	for fi, mvs := range *r.teamPossibleMoves {
		figure := r.bot.team.Figures[fi]
		for _, mv := range *mvs {
			callback := func() bool {
				for _, pos := range *figure.GetBrokenFields() {
					if r.bot.team.Figures.ExistsByPosition(pos) &&
						!r.posIsProtected(pos) &&
						r.posIsAttacked(pos) {
						f := r.bot.team.Figures.GetByPosition(pos)
						rat := r.getFigureValueByName(f.GetName())
						//fmt.Printf("protect figure rat=%f\n\n", rat)
						mv.SetRating(mv.GetRating() + rat)
					}
				}
				return false
			}
			figure.SimulationMove(mv.Position, callback)
		}
	}
}

func (r *rating) getFigureValueByName(figureName string) float64 {
	switch figureName {
	case "pawn":
		return 1
	case "knight":
		return 2
	case "bishop":
		return 2
	case "rook":
		return 3
	case "queen":
		return 4
	}
	return 0
}

// figureIsProtected returns true if the position is protected otherwise returns false
func (r *rating) posIsProtected(pos *game.Position) bool {
	for _, tbfs := range *r.teamBrokenFields {
		for _, tbf := range *tbfs {
			if *pos == *tbf {
				return true
			}
		}
	}
	return false
}

// posIsAttacked returns true if the position is attacked otherwise returns false
func (r *rating) posIsAttacked(pos *game.Position) bool {
	for _, ebfs := range *r.enemyBrokenFields {
		for _, ebf := range *ebfs {
			if *pos == *ebf {
				return true
			}
		}
	}
	return false
}

// GetMoveWithMaxRating get a link to the new action with the highest rating
func (r *rating) getMoveWithMaxRating() *move {
	var rat float64 = -99999
	var m *game.Move
	var fi game.FigurerIndex
	for index, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			if mv.GetRating() > rat {
				rat = mv.GetRating()
				m = mv
				fi = index
			}
		}
	}
	return newMove(r.bot, r.bot.team.Figures[fi].GetPosition(), game.NewPosition(m.X, m.Y))
}
