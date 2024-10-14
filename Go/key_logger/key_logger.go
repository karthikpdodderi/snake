package key_logger

type KeyLogger interface {
	Start()
	Stop()
	Get() rune
}
