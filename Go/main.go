package main

import (
	"board"
	"logic"
	"time"
)

func main() {
	logicInterface := logic.NewLogic(10*time.Millisecond, 20, 10, 10, 5, board.LEFT, '*', 'M', '.', 'q', 'w', 's', 'a', 'd', 300*time.Millisecond, false)
	logicInterface.Start()
}
