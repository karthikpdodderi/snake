package main

import (
	"board"
	"logic"
	"time"
)

func main() {
	logicInterface := logic.NewLogic(10*time.Millisecond, 20, 10, 20, 5, board.LEFT, '#', 'M', '.', 'q', 'w', 's', 'a', 'd', 'p', 150*time.Millisecond, false)
	logicInterface.Start()
}
