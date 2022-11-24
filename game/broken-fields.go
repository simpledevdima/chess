package game

import (
	"errors"
	"fmt"
	"log"
)

type BrokenFieldIndex int

// BrokenFields map referencing positions
type BrokenFields map[BrokenFieldIndex]*Position

// Clear remove all positions from the map
func (p *BrokenFields) Clear() {
	for index := range *p {
		delete(*p, index)
	}
}

// Set sets the position to the map at the specified index
func (p *BrokenFields) Set(index BrokenFieldIndex, pos *Position) BrokenFieldIndex {
	(*p)[index] = pos
	index++
	return index
}

// Get returns the position by index in map
func (p *BrokenFields) Get(index BrokenFieldIndex) *Position {
	pos := (*p)[index]
	if pos == nil {
		log.Println(errors.New(fmt.Sprintf("position with index %d not found", index)))
	}
	return pos
}

// Remove delete link to the position from the map at specified index
func (p *BrokenFields) Remove(index BrokenFieldIndex) {
	if (*p)[index] == nil {
		log.Println(errors.New(fmt.Sprintf("position with index %d not found", index)))
	} else {
		delete(*p, index)
	}
}
