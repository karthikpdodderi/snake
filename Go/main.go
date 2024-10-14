package main

import (
	"board"
	"logic"
	"time"
)

func main() {
	logicInterface := logic.NewLogic(10*time.Millisecond, 20, 4, 4, 5, board.LEFT, '*', 'M', '.', 'q', 'w', 's', 'a', 'd', 1000*time.Millisecond, false)
	logicInterface.Start()
}
