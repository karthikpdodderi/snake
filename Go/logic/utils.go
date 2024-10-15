package logic

import (
	"board"
	"fmt"
)

func (data *logicData) display() {
	data.printer.Print()
	defer func() {
		data.printer.Clear()
		fmt.Printf("Game state : %s \n", data.stateKeeper.GetState().ToString())
		fmt.Println("Toal mice caught : ", data.miceCounter.GetMiceCount())
		fmt.Println("Press any key to exit ... ")
	}()
	for {
		select {
		case <-*data.refreshChan:
			data.printer.Clear()
			data.printer.Print()
		case <-*data.clearChan:
			return
		}
	}
}

func (data *logicData) getDirection(keyPressed rune) board.Direction {
	switch keyPressed {
	case data.upRune:
		return board.UP
	case data.downRune:
		return board.DOWN
	case data.leftRune:
		return board.LEFT
	case data.rightRune:
		return board.RIGHT
	default:
		panic(fmt.Sprintf("invalid direction for key pressed %v ", keyPressed))
	}
}
