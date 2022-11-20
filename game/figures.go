package game

import (
	"errors"
	"fmt"
	"log"
)

type FigureIndex int

// Figures interface map referencing figures
type Figures map[FigureIndex]Figure

// Clear remove all Figures from the map
func (f *Figures) Clear() {
	for index := range *f {
		delete(*f, index)
	}
}

// Set sets the figure's interface to the map at the specified index
func (f *Figures) Set(index FigureIndex, figure Figure) {
	(*f)[index] = figure
}

// Get returns the figure interface by index in map
func (f *Figures) Get(index FigureIndex) Figure {
	figure := (*f)[index]
	if figure == nil {
		log.Println(errors.New(fmt.Sprintf("figure with index %d not found", index)))
	}
	return figure
}

// RemoveByIndex remove figure interface from map at specified index
func (f *Figures) RemoveByIndex(index FigureIndex) {
	if (*f)[index] == nil {
		log.Println(errors.New(fmt.Sprintf("figure with index %d not found", index)))
	} else {
		delete(*f, index)
	}
}

// GetIndexAndFigureByCoords returns the index and the figure interface of the figure interface at the specified coordinates
func (f *Figures) GetIndexAndFigureByCoords(x, y int) (FigureIndex, Figure) {
	for index, figure := range *f {
		fx, fy := figure.GetPosition().Get()
		if fx == x && fy == y {
			return index, figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", x, y)))
	return 0, nil
}

// GetByCoords returns the figure interface at the specified coordinates
func (f *Figures) GetByCoords(x, y int) Figure {
	for _, figure := range *f {
		fx, fy := figure.GetPosition().Get()
		if fx == x && fy == y {
			return figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", x, y)))
	return nil
}

// GetByName returns the first figure interface found by shape name
func (f *Figures) GetByName(name string) Figure {
	for _, figure := range *f {
		if name == figure.GetName() {
			return figure
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with name %s not found", name)))
	return nil
}

// GetIndexByName returns the first figure interface index found by figure name
func (f *Figures) GetIndexByName(name string) FigureIndex {
	for index, figure := range *f {
		if name == figure.GetName() {
			return index
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with name %s not found", name)))
	return 0
}

// GetIndexByCoords returns the index of the figure interface at the specified coordinates
func (f *Figures) GetIndexByCoords(x, y int) FigureIndex {
	for index, figure := range *f {
		fx, fy := figure.GetPosition().Get()
		if fx == x && fy == y {
			return index
		}
	}
	log.Println(errors.New(fmt.Sprintf("figure with coordinates %dx%d not found", x, y)))
	return 0
}

// ExistsByCoords returns true if the shape interface exists in the map at the specified coordinates, otherwise returns false
func (f *Figures) ExistsByCoords(x int, y int) bool {
	for _, figure := range *f {
		figX, figY := figure.GetPosition().Get()
		if figX == x && figY == y {
			return true
		}
	}
	return false
}
