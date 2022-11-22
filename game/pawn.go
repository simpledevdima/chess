package game

func NewPawn(pos *Position, t *Team) *Pawn {
	p := &Pawn{}
	p.figurer = p
	p.Position = pos
	p.SetName("pawn")
	p.SetTeam(t)
	return p
}

// Pawn is data type of chess figure
type Pawn struct {
	Figure
}

// GetBrokenFields return a slice of Positions with broken fields
func (p *Pawn) GetBrokenFields() *Positions {
	poss := make(Positions)
	var pi PositionIndex
	switch p.team.Name {
	case White:
		pos := NewPosition(p.X+1, p.Y+1)
		if p.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		pos = NewPosition(p.X-1, p.Y+1)
		if p.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
	case Black:
		pos := NewPosition(p.X+1, p.Y-1)
		if p.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
		pos = NewPosition(p.X-1, p.Y-1)
		if p.positionOnBoard(pos) {
			pi = poss.Set(pi, pos)
		}
	}
	return &poss
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (p *Pawn) GetPossibleMoves(thereIs bool) *Positions {
	poss := make(Positions)
	var pi PositionIndex
	switch p.team.Name {
	case White:
		pos1 := NewPosition(p.X, p.Y+1)
		if p.positionOnBoard(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) {
			pi = poss.Set(pi, pos1)
			if thereIs {
				return &poss
			}
		}
		pos2 := NewPosition(p.X, p.Y+2)
		if p.positionOnBoard(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) {
			pi = poss.Set(pi, pos2)
			if thereIs {
				return &poss
			}
		}
	case Black:
		pos1 := NewPosition(p.X, p.Y-1)
		if p.positionOnBoard(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) {
			pi = poss.Set(pi, pos1)
			if thereIs {
				return &poss
			}
		}
		pos2 := NewPosition(p.X, p.Y-2)
		if p.positionOnBoard(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) {
			pi = poss.Set(pi, pos2)
			if thereIs {
				return &poss
			}
		}
	}
	for _, position := range *p.GetBrokenFields() {
		if (p.team.enemy.Figures.ExistsByPosition(position) ||
			p.team.enemy.pawnDoubleMove.isTakeOnThePass(position)) &&
			!p.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
			if thereIs {
				return &poss
			}
		}
	}
	return &poss
}

// CanWalkLikeThat desc
func (p *Pawn) CanWalkLikeThat(pos *Position) bool {
	switch p.team.Name {
	case White:
		switch {
		case p.X == pos.X && (p.Y+1 == pos.Y || p.Y+2 == pos.Y):
			return true // right move
		case (p.X == pos.X+1 || p.X == pos.X-1) && p.Y+1 == pos.Y && p.team.enemy.Figures.ExistsByPosition(pos):
			return true // right eating
		}
	case Black:
		switch {
		case p.X == pos.X && (p.Y-1 == pos.Y || p.Y-2 == pos.Y):
			return true // right move
		case (p.X == pos.X+1 || p.X == pos.X-1) && p.Y-1 == pos.Y && p.team.enemy.Figures.ExistsByPosition(pos):
			return true // right eating
		}
	}
	return false
}
