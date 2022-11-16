package game

func NewPosition(x, y int) *Position {
	p := &Position{}
	p.Set(x, y)
	return p
}

// Position data type containing the coordinates on the board
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Set position to argument values
func (p *Position) Set(x, y int) {
	p.X = x
	p.Y = y
}

// Get return position
func (p *Position) Get() (int, int) {
	return p.X, p.Y
}
