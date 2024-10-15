package logic

import (
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
