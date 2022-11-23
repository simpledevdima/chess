package game

// NewMove возвращает ссылку на новый ход с указаной ссылкой на переданную в аргументе позицию
func NewMove(pos *Position) *Move {
	m := &Move{}
	m.Position = pos
	return m
}

// Move тип данных содержащий ссылку на позицию хода и поле рейтинга
type Move struct {
	*Position
	Rating float64 `json:"-"`
}

// SetRating sets the rating value from the argument
func (m *Move) SetRating(rat float64) {
	m.Rating = rat
}

// GetRating gets the rating value
func (m *Move) GetRating() float64 {
	return m.Rating
}
