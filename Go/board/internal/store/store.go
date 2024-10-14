package store

import "board/internal/position"

type Store interface {
	Add(position.Position) error
	Remove(position.Position) error
	Random() (position.Position, error)
}
