package logger

type Logger interface {
	Print(msg string)
	Close()
}
