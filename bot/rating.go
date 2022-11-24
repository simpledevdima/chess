package chess_bot

import (
	"github.com/skvdmt/chess/game"
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

// EatUnprotectedFigure метод изменяющий рейтинг хода на основе того, что фигура есть незащищенную фигуру противника
func (r *rating) EatUnprotectedFigure() {
	for _, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			for _, ef := range r.bot.enemy.Figures {
				if *mv.Position == *ef.GetPosition() {
					var protected bool
					func() {
						for _, ebfs := range *r.enemyBrokenFields {
							for _, ebf := range *ebfs {
								if *ebf == *mv.Position {
									protected = true
									return
								}
							}
						}
					}()
					if !protected {
						var rat float64
						switch ef.GetName() {
						case "pawn":
							rat = 1
						case "knight":
							rat = 2
						case "bishop":
							rat = 2
						case "rook":
							rat = 3
						case "queen":
							rat = 4
						}
						mv.SetRating(mv.GetRating() + rat)
					}
				}
			}
		}
	}
}

// MoveToBrokenField changing the rating for moves to broken opponent's fields
func (r *rating) MoveToBrokenField() {
	for fi, mvs := range *r.teamPossibleMoves {
		for _, mv := range *mvs {
			var broken bool
			func() {
				for _, ebfs := range *r.enemyBrokenFields {
					for _, ebf := range *ebfs {
						if *ebf == *mv.Position {
							broken = true
							return
						}
					}
				}
			}()
			if broken {
				var rat float64
				switch r.bot.team.Figures[fi].GetName() {
				case "pawn":
					rat = 1
				case "knight":
					rat = 2
				case "bishop":
					rat = 2
				case "rook":
					rat = 3
				case "queen":
					rat = 4
				}
				mv.SetRating(mv.GetRating() - rat)
			}

		}
	}
}

// remove an unprotected piece from the opponent's beaten square
func (r *rating) moveUnprotectedFigureFromBrokenFieldToProtectedOrSecureField() {

}

// figureIsProtected returns true if the position is protected otherwise returns false
func (r *rating) posIsProtected(pos *game.Position) bool {
	var protected bool
	func() {
		for _, tbfs := range *r.teamBrokenFields {
			for _, tbf := range *tbfs {
				if *pos == *tbf {
					protected = true
					return
				}
			}
		}
	}()
	return protected
}

// posIsAttacked returns true if the position is attacked otherwise returns false
func (r *rating) posIsAttacked(pos *game.Position) bool {
	var attacking bool
	func() {
		for _, ebfs := range *r.enemyBrokenFields {
			for _, ebf := range *ebfs {
				if *pos == *ebf {
					attacking = true
					return
				}
			}
		}
	}()
	return attacking
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
