package logic

import (
	"board"
	"fmt"
	"key_logger"
	"logger"
	"sync"
	"time"
)

type logicData struct {
	quitRune           rune
	printer            board.Displayer
	mover              board.Mover
	miceCounter        board.ScoreKeeper
	stateKeeper        board.StateKeeper
	keyLogger          key_logger.KeyLogger
	upRune             rune
	downRune           rune
	leftRune           rune
	rightRune          rune
	pausePlayRune      rune
	refreshChan        *chan bool
	clearChan          *chan bool
	quitChan           *chan bool
	pausePlayChan      *chan bool
	fileLogger         logger.Logger
	snakeMoveDelayTime time.Duration
	stateMux           sync.Mutex
}

func NewLogic(keyLoggerDelay time.Duration, keyLoggerBuffer int, numRows int, numCols int, snakeInitialLenght int, snakeInitialDir board.Direction, snakeRune rune, mouseRune rune, fieldRune rune, quitRune rune, upRune rune, downRune rune, leftRune rune, rightRune rune, pausePlayRune rune, snakeMoveDelayTime time.Duration, isLogRequired bool) Logic {

	if snakeInitialLenght < 2 {
		panic("Expected snake initial lenght >= 2")
	}

	if numRows*numCols <= snakeInitialLenght {
		panic("snake lenght must be less than arena area")
	}

	fileLogName := fmt.Sprintf("%v", time.Now().UnixMilli())
	fileLogger, err := logger.NewFileLogger(fileLogName, isLogRequired)
	if err != nil {
		panic(fmt.Sprintf("Error while bringing up a file logger. Error : %v ", err))
	}

	keyLoggerInterface := key_logger.NewKeyLogger(keyLoggerDelay, keyLoggerBuffer, fileLogger)

	boardMover, boardDisplayer, boardScoreKeeper, stateKeeper := board.NewBoard(numRows, numCols, snakeInitialLenght, snakeInitialDir, snakeRune, mouseRune, fieldRune, fileLogger)

	refreshChan := make(chan bool, 1)
	quitChan := make(chan bool, 1)
	clearChan := make(chan bool, 1)
	pausePlayChan := make(chan bool, 1)

	return &logicData{
		quitRune:           quitRune,
		printer:            boardDisplayer,
		mover:              boardMover,
		miceCounter:        boardScoreKeeper,
		stateKeeper:        stateKeeper,
		keyLogger:          keyLoggerInterface,
		refreshChan:        &refreshChan,
		quitChan:           &quitChan,
		clearChan:          &clearChan,
		pausePlayChan:      &pausePlayChan,
		upRune:             upRune,
		downRune:           downRune,
		leftRune:           leftRune,
		rightRune:          rightRune,
		pausePlayRune:      pausePlayRune,
		fileLogger:         fileLogger,
		snakeMoveDelayTime: snakeMoveDelayTime,
		stateMux:           sync.Mutex{},
	}

}

func (data *logicData) Start() {
	data.keyLogger.Start()

	defer data.keyLogger.Stop()
	defer data.fileLogger.Close()

	go data.display()

	go func() {
		isPaused := false
		for {
			select {
			case <-*data.pausePlayChan:
				isPaused = !isPaused
			default:
				if isPaused {
					<-*data.pausePlayChan
					isPaused = !isPaused
				}
				state := data.mover.Continue()
				if state == board.LOSE || state == board.WIN {
					*data.clearChan <- true
					*data.quitChan <- true
					return
				}
				*data.refreshChan <- true
				time.Sleep(data.snakeMoveDelayTime)
			}
		}
	}()

	isPaused := false
	for {
		select {
		case <-*data.quitChan:
			return
		default:
			keyPressed := data.keyLogger.Get()
			switch keyPressed {
			case data.upRune, data.downRune, data.leftRune, data.rightRune:
				if isPaused {
					continue
				}
				data.mover.Turn(data.getDirection(keyPressed))
				*data.refreshChan <- true
			case data.quitRune:
				*data.clearChan <- true
				return
			case data.pausePlayRune:
				isPaused = !isPaused
				*data.pausePlayChan <- true
			}
		}
	}

}
