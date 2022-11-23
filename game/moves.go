package game

import (
	"errors"
	"fmt"
	"log"
)

type MoveIndex int

// Moves map referencing moves
type Moves map[MoveIndex]*Move

// Clear remove all Moves from the map
func (m *Moves) Clear() {
	for index := range *m {
		delete(*m, index)
	}
}

// Set sets the move to the map at the specified index
func (m *Moves) Set(index MoveIndex, mov *Move) MoveIndex {
	(*m)[index] = mov
	index++
	return index
}

// Get returns the move by index in map
func (m *Moves) Get(index MoveIndex) *Move {
	mov := (*m)[index]
	if mov == nil {
		log.Println(errors.New(fmt.Sprintf("move with index %d not found", index)))
	}
	return mov
}

// Remove delete link to the move from the map at specified index
func (m *Moves) Remove(index MoveIndex) {
	if (*m)[index] == nil {
		log.Println(errors.New(fmt.Sprintf("move with index %d not found", index)))
	} else {
		delete(*m, index)
	}
}
