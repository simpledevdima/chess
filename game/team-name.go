package game

// TeamName command data type
type TeamName int

const (
	White TeamName = iota
	Black
	Spectators
)

// string return team name as string
func (t *TeamName) String() string {
	return [...]string{"white", "black", "spectators"}[*t]
}
