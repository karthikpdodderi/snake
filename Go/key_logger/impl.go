package key_logger

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"logger"
	"time"
)

type keyLoggerData struct {
	waitTime   time.Duration
	stopper    chan bool
	buffer     chan rune
	fileLogger logger.Logger
}

func NewKeyLogger(waitTime time.Duration, loggerBufferLenght int, fileLogger logger.Logger) KeyLogger {
	return &keyLoggerData{
		waitTime:   waitTime,
		stopper:    make(chan bool),
		buffer:     make(chan rune, loggerBufferLenght),
		fileLogger: fileLogger,
	}
}

func (data *keyLoggerData) Start() {
	data.fileLogger.Print("opening keyboard")
	err := keyboard.Open()
	if err != nil {
		panic(fmt.Sprintf("Error while opening the keyboard. Error : %v ", err))
	}
	go func() {
		for {
			select {
			case <-data.stopper:
				// data.fileLogger.Print("closing keyboard")
				// err = keyboard.Close()
				// if err != nil {
				// 	panic(fmt.Sprintf("Error while closing the keyboard. Error : %v ", err))
				// }
				return
			default:
				char, _, err := keyboard.GetSingleKey()
				if err != nil {
					panic(fmt.Sprintf("Error while getting single key from keyboard. Error : %v ", err))
				}
				data.buffer <- char
			}
			time.Sleep(data.waitTime)
			// waiting for process, i.e the charecter pressed
			// result in stopping the channel
		}
	}()
}

func (data *keyLoggerData) Stop() {
	data.stopper <- true
}

func (data *keyLoggerData) Get() rune {
	data.fileLogger.Print("waiting inside get ...")
	val := <-data.buffer
	data.fileLogger.Print(fmt.Sprintf("inside get got %s ", string(val)))
	return val
}
