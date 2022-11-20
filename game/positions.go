package game

import (
	"errors"
	"fmt"
	"log"
)

type PositionIndex int

// Positions map referencing positions
type Positions map[PositionIndex]*Position

// Clear remove all positions from the map
func (p *Positions) Clear() {
	for index := range *p {
		delete(*p, index)
	}
}

// Set sets the position to the map at the specified index
func (p *Positions) Set(index PositionIndex, pos *Position) PositionIndex {
	(*p)[index] = pos
	index++
	return index
}

// Get returns the position by index in map
func (p *Positions) Get(index PositionIndex) *Position {
	pos := (*p)[index]
	if pos == nil {
		log.Println(errors.New(fmt.Sprintf("position with index %d not found", index)))
	}
	return pos
}

// Remove delete link to the position from the map at specified index
func (p *Positions) Remove(index PositionIndex) {
	if (*p)[index] == nil {
		log.Println(errors.New(fmt.Sprintf("position with index %d not found", index)))
	} else {
		delete(*p, index)
	}
}
