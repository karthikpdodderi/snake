package queue

import (
	"board/internal/position"
	"fmt"
	"logger"
)

type queueData struct {
	isPositionPresent map[position.Position]bool
	queue             []position.Position
	fileLogger        logger.Logger
}

func NewQueue(fileLogger logger.Logger) Queue {
	queue := &queueData{
		isPositionPresent: make(map[position.Position]bool),
		queue:             make([]position.Position, 0),
		fileLogger:        fileLogger,
	}
	return queue
}

func (data *queueData) Add(pos position.Position) (err error) {
	_, isPresent := data.isPositionPresent[pos]
	if isPresent {
		return fmt.Errorf("position %v already present in the queue", pos)
	}
	data.queue = append(data.queue, pos)
	return nil
}

func (data *queueData) Remove() (position.Position, error) {
	if len(data.queue) == 0 {
		return position.Position{}, fmt.Errorf("No data left in queue to remove")
	}
	pos := data.queue[0]
	delete(data.isPositionPresent, pos)
	data.queue = data.queue[1:]
	return pos, nil
}

func (data *queueData) GetLength() int {
	return len(data.queue)
}

func (data *queueData) GetHead() (position.Position, error) {
	if len(data.queue) == 0 {
		return position.Position{}, fmt.Errorf("No data left present in queue")
	}
	return data.queue[len(data.queue)-1], nil
}
