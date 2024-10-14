package board

import (
	"board/internal/position"
	"board/internal/queue"
	"board/internal/store"
	"fmt"
	"logger"
	"sync"
)

type boardData struct {
	arenaMux  sync.Mutex
	arenaData arenaData
}

type arenaData struct {
	arena             [][]rune
	snakeHeadDir      Direction
	snake             queue.Queue
	fieldBuffer       store.Store
	snakeRune         rune
	mouseRune         rune
	fieldRune         rune
	numRows           int
	numColumns        int
	miceCaught        int
	fileLogger        logger.Logger
	isSnakeDirChanged bool
}

func NewBoard(numRows int, numColumns int, snakeInitialLenght int, snakeIntialDirection Direction, snakeRune rune, mouseRune rune, fieldRune rune, fileLogger logger.Logger) (Mover, Displayer, ScoreKeeper) {

	arena := make([][]rune, 0)
	fieldStore := store.NewStore(fileLogger)

	for i := 0; i < numRows; i++ {
		row := make([]rune, numColumns)
		for j := range row {
			row[j] = fieldRune
			fieldStore.Add(position.Position{RowNum: i, ColNum: j})
		}
		arena = append(arena, row)
	}

	snakeInitialPos := position.Position{
		RowNum: numRows / 2,
		ColNum: numColumns / 2,
	}

	snakeEndingPos := snakeInitialPos

	switch snakeIntialDirection {
	case UP:
		endRow := snakeEndingPos.RowNum + (snakeInitialLenght - 1)
		if endRow >= numRows {
			endRow = numRows - 1
		}
		snakeEndingPos.RowNum = endRow

	case DOWN:
		startRow := snakeEndingPos.RowNum - (snakeInitialLenght - 1)
		if startRow < 0 {
			startRow = 0
		}
		snakeInitialPos.RowNum = startRow

	case LEFT:
		endColumn := snakeEndingPos.ColNum + (snakeInitialLenght - 1)
		if endColumn >= numColumns {
			endColumn = numColumns - 1
		}
		snakeEndingPos.ColNum = endColumn

	case RIGHT:
		startColumn := snakeEndingPos.ColNum - (snakeInitialLenght - 1)
		if startColumn < 0 {
			startColumn = 0
		}
		snakeInitialPos.ColNum = startColumn
	}

	snakeQueue := queue.NewQueue(fileLogger)
	switch snakeIntialDirection {
	case UP, DOWN:
		colNumber := snakeInitialPos.ColNum
		for i := snakeEndingPos.RowNum; i >= snakeInitialPos.RowNum; i-- {
			pos := position.Position{RowNum: i, ColNum: colNumber}
			arena[i][colNumber] = snakeRune
			snakeQueue.Add(pos)
			err := fieldStore.Remove(pos)
			if err != nil {
				panic(fmt.Sprintf("Error while removing pos %v form field store. Error : %v ", pos, err))
			}
		}
	case LEFT, RIGHT:
		rowNumber := snakeInitialPos.RowNum
		for j := snakeEndingPos.ColNum; j >= snakeInitialPos.ColNum; j-- {
			pos := position.Position{RowNum: rowNumber, ColNum: j}
			arena[rowNumber][j] = snakeRune
			snakeQueue.Add(pos)
			err := fieldStore.Remove(pos)
			if err != nil {
				panic(fmt.Sprintf("Error while removing pos %v form field store. Error : %v ", pos, err))
			}
		}
	}

	mousePos, err := fieldStore.Random()
	if err != nil {
		panic(fmt.Sprintf("Error in getting random position for mouse. Error : %v ", err))
	}
	err = fieldStore.Remove(mousePos)
	if err != nil {
		panic(fmt.Sprintf("Error while removing mouse pos %v form field store. Error : %v ", mousePos, err))
	}
	arena[mousePos.RowNum][mousePos.ColNum] = mouseRune

	boardInterface := boardData{
		arenaMux: sync.Mutex{},
		arenaData: arenaData{
			arena:             arena,
			snakeHeadDir:      snakeIntialDirection,
			snake:             snakeQueue,
			numRows:           numRows,
			numColumns:        numColumns,
			miceCaught:        0,
			fileLogger:        fileLogger,
			fieldBuffer:       fieldStore,
			snakeRune:         snakeRune,
			mouseRune:         mouseRune,
			fieldRune:         fieldRune,
			isSnakeDirChanged: false,
		},
	}
	return &boardInterface, &boardInterface, &boardInterface
}

func (data *boardData) GetMiceCount() int {

	data.arenaMux.Lock()
	defer data.arenaMux.Unlock()

	return data.arenaData.miceCaught
}

func (data *boardData) Continue() State {

	data.arenaMux.Lock()
	defer data.arenaMux.Unlock()
	defer func() {
		data.arenaData.isSnakeDirChanged = false
	}()

	return data.carryOn()
}

func (data *boardData) Turn(dir Direction) {

	data.arenaMux.Lock()
	defer data.arenaMux.Unlock()

	if !data.arenaData.isSnakeDirChanged {
		switch data.arenaData.snakeHeadDir {
		case UP, DOWN:
			if dir == LEFT || dir == RIGHT {
				data.arenaData.snakeHeadDir = dir
				data.arenaData.isSnakeDirChanged = true
			}
		case LEFT, RIGHT:
			if dir == UP || dir == DOWN {
				data.arenaData.snakeHeadDir = dir
				data.arenaData.isSnakeDirChanged = true
			}
		}
	}
}

func (data *boardData) Print() {

	data.arenaMux.Lock()
	defer data.arenaMux.Unlock()

	snakeHeadPos, err := data.arenaData.snake.GetHead()
	if err != nil {
		panic(fmt.Sprintf("Error while getting head element. Error : %v ", err))
	}

	for rowNum, row := range data.arenaData.arena {
		for colNum, cell := range row {
			if rowNum == snakeHeadPos.RowNum && colNum == snakeHeadPos.ColNum {
				switch data.arenaData.snakeHeadDir {
				case UP:
					fmt.Printf("%v ", "^")
				case DOWN:
					fmt.Printf("%v ", "v")
				case LEFT:
					fmt.Printf("%v ", "<")
				case RIGHT:
					fmt.Printf("%v ", ">")
				default:
					panic(fmt.Sprintf("invalid snake head dir %v ", data.arenaData.snakeHeadDir))
				}
			} else {
				fmt.Printf("%v ", string(cell))
			}
		}
		fmt.Println()
	}

}

func (data *boardData) Clear() {

	data.arenaMux.Lock()
	defer data.arenaMux.Unlock()

	for i := 0; i < len(data.arenaData.arena); i++ {
		fmt.Print("\033[A\033[K") // Move up and clear the line
	}

}
