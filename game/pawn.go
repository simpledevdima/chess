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
	figureData
}

// DetectionOfPossibleMove return slice of Position with coords for possible moves
func (p *Pawn) DetectionOfPossibleMove() []*Position {
	var data []*Position
	switch p.team.Name {
	case White:
		pos1 := NewPosition(p.X, p.Y+1)
		//fmt.Println(p.positionOnBoard(pos1))
		//fmt.Println(p.kingOnTheBeatenFieldAfterMove(pos1))
		//fmt.Println(p.team.Figures.ExistsByPosition(pos1))
		//fmt.Println(p.team.enemy.Figures.ExistsByPosition(pos1))
		if p.positionOnBoard(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) {
			data = append(data, pos1)
		}
		pos2 := NewPosition(p.X, p.Y+2)
		if p.positionOnBoard(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) {
			data = append(data, pos2)
		}
	case Black:
		pos1 := NewPosition(p.X, p.Y-1)
		if p.positionOnBoard(pos1) &&
			!p.kingOnTheBeatenFieldAfterMove(pos1) &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) {
			data = append(data, pos1)
		}
		pos2 := NewPosition(p.X, p.Y-2)
		if p.positionOnBoard(pos2) &&
			!p.kingOnTheBeatenFieldAfterMove(pos2) &&
			!p.IsAlreadyMove() &&
			!p.team.Figures.ExistsByPosition(pos1) &&
			!p.team.Figures.ExistsByPosition(pos2) &&
			!p.team.enemy.Figures.ExistsByPosition(pos1) &&
			!p.team.enemy.Figures.ExistsByPosition(pos2) {
			data = append(data, pos2)
		}
	}
	for _, position := range p.detectionOfBrokenFields() {
		if (p.team.enemy.Figures.ExistsByPosition(position) ||
			p.team.enemy.pawnDoubleMove.isTakeOnThePass(position)) &&
			!p.kingOnTheBeatenFieldAfterMove(position) {
			data = append(data, position)
		}
	}
	return data
}

// detectionOfBrokenFields return a slice of Positions with broken fields
func (p *Pawn) detectionOfBrokenFields() []*Position {
	var data []*Position
	switch p.team.Name {
	case White:
		pos := NewPosition(p.X+1, p.Y+1)
		if p.positionOnBoard(pos) {
			data = append(data, pos)
		}
		pos = NewPosition(p.X-1, p.Y+1)
		if p.positionOnBoard(pos) {
			data = append(data, pos)
		}
	case Black:
		pos := NewPosition(p.X+1, p.Y-1)
		if p.positionOnBoard(pos) {
			data = append(data, pos)
		}
		pos = NewPosition(p.X-1, p.Y-1)
		if p.positionOnBoard(pos) {
			data = append(data, pos)
		}
	}
	return data
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
	for _, position := range p.detectionOfBrokenFields() {
		if *position == *pos &&
			(p.team.enemy.Figures.ExistsByPosition(pos) || p.team.enemy.pawnDoubleMove.isTakeOnThePass(pos)) {
			return true, ""
		}
	}
	// move pawn
	for _, position := range p.DetectionOfPossibleMove() {
		if *position == *pos {
			return true, ""
		}
	}
	return false, "this figure cant make that move"
}

// Move change Position of figure to Position from arguments
func (p *Pawn) Move(pos *Position) {
	p.team.enemy.pawnDoubleMove.pawnTakeOnThePass(pos)
	p.team.pawnDoubleMove.clearPawnDoubleMove()
	p.team.pawnDoubleMove.pawnMakesDoubleMove(p, p.GetPosition(), pos)
	p.MoveFigure(pos)
	p.transformPawnTOQueen(pos)
}

// transformPawnToQueen promote a pawn to a queen
func (p *Pawn) transformPawnTOQueen(pos *Position) {
	if p.Y == 1 || p.Y == 8 {
		figureID := p.team.Figures.GetIndexByPosition(pos)
		// replace pawn to queen
		p.team.Figures.Set(figureID, NewQueen(pos, p.team))
		p.team.Figures.Get(figureID).setAlreadyMove(true)
	}
}
