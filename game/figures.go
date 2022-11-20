package game

import (
	"errors"
	"fmt"
	"log"
)

type FigurerIndex int

// Figures interface map referencing figures
type Figures map[FigurerIndex]Figurer

// Clear remove all Figures from the map
func (f *Figures) Clear() {
	for index := range *f {
		delete(*f, index)
	}
}

// Set sets the figure's interface to the map at the specified index
func (f *Figures) Set(index FigurerIndex, figure Figurer) {
	(*f)[index] = figure
}

// Get returns the figure interface by index in map
func (f *Figures) Get(index FigurerIndex) Figurer {
	figure := (*f)[index]
	if figure == nil {
		log.Println(errors.New(fmt.Sprintf("figure with index %d not found", index)))
	}
	return figure
}

// RemoveByIndex remove figure interface from map at specified index
func (f *Figures) RemoveByIndex(index FigurerIndex) {
	if (*f)[index] == nil {
		log.Println(errors.New(fmt.Sprintf("figure with index %d not found", index)))
	} else {
		delete(*f, index)
	}
}

// GetIndexAndFigureByPosition returns the index and the figure interface of the figure interface at the specified coordinates
func (f *Figures) GetIndexAndFigureByPosition(pos *Position) (FigurerIndex, Figurer) {
	for index, figure := range *f {
		if *figure.GetPosition() == *pos {
			return index, figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", pos.X, pos.Y)))
	return 0, nil
}

// GetByPosition returns the figure interface at the specified coordinates
func (f *Figures) GetByPosition(pos *Position) Figurer {
	for _, figure := range *f {
		figPos := figure.GetPosition()
		if *figPos == *pos {
			return figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", pos.X, pos.Y)))
	return nil
}

// GetByName returns the first figure interface found by shape name
func (f *Figures) GetByName(name string) Figurer {
	for _, figure := range *f {
		if name == figure.GetName() {
			return figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with name %s not found", name)))
	return nil
}

// GetIndexByName returns the first figure interface index found by figure name
func (f *Figures) GetIndexByName(name string) FigurerIndex {
	for index, figure := range *f {
		if name == figure.GetName() {
			return index
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with name %s not found", name)))
	return 0
}

// GetIndexByPosition returns the index of the figure interface at the specified coordinates
func (f *Figures) GetIndexByPosition(pos *Position) FigurerIndex {
	for index, figure := range *f {
		figPos := figure.GetPosition()
		if *figPos == *pos {
			return index
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", pos.X, pos.Y)))
	return 0
}

// ExistsByPosition returns true if the shape interface exists in the map at the specified coordinates, otherwise returns false
func (f *Figures) ExistsByPosition(pos *Position) bool {
	for _, figure := range *f {
		if *figure.GetPosition() == *pos {
			return true
		}
	}
	return false
}
