package store

import (
	"board/internal/position"
	"fmt"
	"logger"
	"math/rand"
)

type storeData struct {
	positions  []position.Position
	posToIndex map[position.Position]int
	fileLogger logger.Logger
}

func NewStore(fileLogger logger.Logger) Store {
	return &storeData{
		positions:  make([]position.Position, 0),
		posToIndex: make(map[position.Position]int),
		fileLogger: fileLogger,
	}
}

func (data *storeData) Add(pos position.Position) (err error) {
	_, isPresent := data.posToIndex[pos]
	if isPresent {
		return fmt.Errorf("position %v already present in the store", pos)
	}
	data.positions = append(data.positions, pos)
	data.posToIndex[pos] = len(data.positions) - 1
	return nil
}

func (data *storeData) Remove(pos position.Position) error {

	index, isPresent := data.posToIndex[pos]
	if !isPresent {
		return fmt.Errorf("position %v is not present in the store", pos)
	}

	// adding last element to current index
	lastIndex := len(data.positions) - 1
	lastElement := data.positions[lastIndex]

	data.positions[index] = lastElement
	data.posToIndex[lastElement] = index

	// deleting last index
	data.positions = data.positions[0:lastIndex]
	delete(data.posToIndex, pos)

	return nil

}

func (data *storeData) Random() (position.Position, error) {
	if len(data.positions) == 0 {
		return position.Position{}, fmt.Errorf("no position present the store")
	}
	index := rand.Intn(len(data.positions))
	return data.positions[index], nil
}
