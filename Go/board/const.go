package board

type State int

const (
	LOSE     State = 0
	WIN      State = 1
	CONTINUE State = 2
)

func (state State) ToString() string {
	switch state {
	case LOSE:
		return "LOSE"
	case WIN:
		return "WIN"
	case CONTINUE:
		return "CONTINUE"
	default:
		return ""
	}
}

type Direction int

const (
	UP    Direction = 0
	DOWN  Direction = 1
	LEFT  Direction = 2
	RIGHT Direction = 3
)
