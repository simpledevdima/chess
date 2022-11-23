package game

// NewPosition returns a link to a new position with the parameters specified in the argument
func NewPosition(x, y uint8) *Position {
	p := &Position{}
	p.Set(x, y)
	return p
}

// Position data type containing the coordinates on the board
type Position struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Set position to argument values
func (p *Position) Set(x, y uint8) {
	p.X = x
	p.Y = y
}

// Get return position
func (p *Position) Get() (uint8, uint8) {
	return p.X, p.Y
}
