package board

import (
	"board/internal/position"
	"fmt"
)

func (data *boardData) addSnakeHead(headPos position.Position) {
	err := data.arenaData.snake.Add(headPos)
	if err != nil {
		panic(fmt.Sprintf("error while adding position %v to queue. Error : %v ", headPos, err))
	}
	// removing only if field rune
	if data.arenaData.arena[headPos.RowNum][headPos.ColNum] == data.arenaData.fieldRune {
		err = data.arenaData.fieldBuffer.Remove(headPos)
		if err != nil {
			panic(fmt.Sprintf("error while removing head position %v from field buffer during head assignment . Error : %v ", headPos, err))
		}
	}
	data.arenaData.arena[headPos.RowNum][headPos.ColNum] = data.arenaData.snakeRune
}

func (data *boardData) removeSnakeTail() {
	tailPos, err := data.arenaData.snake.Remove()
	if err != nil {
		panic(fmt.Sprintf("error while removing from snake queue during tail removal. Error : %v ", err))
	}
	err = data.arenaData.fieldBuffer.Add(tailPos)
	if err != nil {
		panic(fmt.Sprintf("error while adding tail position %v from field buffer during tail removal. Error : %v ", tailPos, err))
	}
	data.arenaData.arena[tailPos.RowNum][tailPos.ColNum] = data.arenaData.fieldRune
}

func (data *boardData) addMouse() {
	mousePos, err := data.arenaData.fieldBuffer.Random()
	if err != nil {
		panic(fmt.Sprintf("error while getting random pos from field buffer during mouse assignment. Error : %v ", err))
	}
	err = data.arenaData.fieldBuffer.Remove(mousePos)
	if err != nil {
		panic(fmt.Sprintf("error while removing mouse position %v from field buffer during mouse assignment . Error : %v ", mousePos, err))
	}
	data.arenaData.arena[mousePos.RowNum][mousePos.ColNum] = data.arenaData.mouseRune
}

func (data *boardData) getNextPos() position.Position {
	nextPos, err := data.arenaData.snake.GetHead()
	if err != nil {
		panic(fmt.Sprintf("Error while getting snake during getting next position. Error : %v ", err))
	}
	switch data.arenaData.snakeHeadDir {
	case UP:
		nextPos.RowNum = (data.arenaData.numRows + (nextPos.RowNum - 1)) % (data.arenaData.numRows)
	case DOWN:
		nextPos.RowNum = (data.arenaData.numRows + (nextPos.RowNum + 1)) % (data.arenaData.numRows)
	case LEFT:
		nextPos.ColNum = (data.arenaData.numColumns + (nextPos.ColNum - 1)) % (data.arenaData.numColumns)
	case RIGHT:
		nextPos.ColNum = (data.arenaData.numColumns + (nextPos.ColNum + 1)) % (data.arenaData.numColumns)
	default:
		panic(fmt.Sprintf("invalid snake head direction %v ", data.arenaData.snakeHeadDir))
	}
	return nextPos
}

func (data *boardData) carryOn() State {

	nextPos := data.getNextPos()
	// if next cell is a snake cell, then declare it a failure
	if data.arenaData.arena[nextPos.RowNum][nextPos.ColNum] == data.arenaData.snakeRune {
		return LOSE
	}

	if data.arenaData.arena[nextPos.RowNum][nextPos.ColNum] == data.arenaData.mouseRune {

		data.addSnakeHead(nextPos)
		data.arenaData.miceCaught++

		// if all cells are occupied by snake, then declare it a victory
		if data.arenaData.snake.GetLength() == (data.arenaData.numColumns * data.arenaData.numRows) {
			return WIN
		}

		data.addMouse()
		return CONTINUE

	}
	data.addSnakeHead(nextPos)
	data.removeSnakeTail()
	return CONTINUE
}
