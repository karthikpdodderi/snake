package board

type Mover interface {
	Continue() State
	Turn(Direction)
}

type Displayer interface {
	Print()
	Clear()
}

type ScoreKeeper interface {
	GetMiceCount() int
}
