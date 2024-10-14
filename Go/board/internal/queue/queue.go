package queue

import "board/internal/position"

type Queue interface {
	Add(position.Position) error
	Remove() (position.Position, error)
	GetLength() int
	GetHead() (position.Position, error)
}
