package logger

import (
	"log"
	"os"
)

type fileLogger struct {
	logger        *log.Logger
	isLogRequired bool
}

func NewFileLogger(filename string, isLogRequired bool) (Logger, error) {

	if !isLogRequired {
		return &fileLogger{logger: nil, isLogRequired: false}, nil
	}
	// Create or open the log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Create a new logger
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &fileLogger{logger: logger, isLogRequired: true}, nil
}

// Log an info message
func (f *fileLogger) Print(message string) {
	if !f.isLogRequired {
		return
	}
	f.logger.Println("[MSG]", message)
}

// Close the logger file
func (f *fileLogger) Close() {
	if !f.isLogRequired {
		return
	}
	if file, ok := f.logger.Writer().(*os.File); ok {
		file.Close()
	}
}
