package game

// Position data type containing the coordinates on the board
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Set position to argument values
func (Position *Position) Set(x, y int) {
	Position.X = x
	Position.Y = y
}

// Get return position
func (Position *Position) Get() (int, int) {
	return Position.X, Position.Y
}
