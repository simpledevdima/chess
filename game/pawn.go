package game

func NewPawn(pos *Position, t *Team) *Pawn {
	p := &Pawn{}
	p.Position = pos
	p.SetName("pawn")
	p.SetTeam(t)
	return p
}

// Pawn is data type of chess figure
type Pawn struct {
	Figure
}

// GetPossibleMoves return slice of Position with coords for possible moves
func (p *Pawn) GetPossibleMoves() *Positions {
	poss := make(Positions)
	var pi PositionIndex
	switch p.team.Name {
	case White:
		pos1 := NewPosition(p.X, p.Y+1)
		if p.positionOnBoard(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) {
			pi = poss.Set(pi, pos1)
		}
		pos2 := NewPosition(p.X, p.Y+2)
		if p.positionOnBoard(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) {
			pi = poss.Set(pi, pos2)
		}
	case Black:
		pos1 := NewPosition(p.X, p.Y-1)
		if p.positionOnBoard(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) {
			pi = poss.Set(pi, pos1)
		}
		pos2 := NewPosition(p.X, p.Y-2)
		if p.positionOnBoard(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) {
			pi = poss.Set(pi, pos2)
		}
	}
	for _, position := range *p.GetBrokenFields() {
		if (p.team.enemy.Figures.ExistsByPosition(position) ||
			p.team.enemy.pawnDoubleMove.isTakeOnThePass(position)) &&
			!p.kingOnTheBeatenFieldAfterMove(position) {
			pi = poss.Set(pi, position)
		}
	}
	return &poss
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

// Validation return true if this move are valid or return false
func (p *Pawn) Validation(pos *Position) (bool, string) {
	if !p.positionOnBoard(pos) {
		return false, "attempt to go out the board"
	}
	if *p.GetPosition() == *pos {
		return false, "can't walk around"
	}
	if p.team.Figures.ExistsByPosition(pos) {
		return false, "this place is occupied by your figure"
	}
	if p.kingOnTheBeatenFieldAfterMove(pos) {
		return false, "your king stands on a beaten field"
	}
	// detect Position for eat and check it for input data eat coords
	for _, position := range *p.GetBrokenFields() {
		if *position == *pos &&
			(p.team.enemy.Figures.ExistsByPosition(pos) || p.team.enemy.pawnDoubleMove.isTakeOnThePass(pos)) {
			return true, ""
		}
	}
	// move pawn
	for _, position := range *p.GetPossibleMoves() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}
