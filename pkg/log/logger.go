package log

import (
	"log"
	"os"
)

type Level byte

const (
	Debug Level = 1
	Error Level = 2
)

type Logger struct {
	debugLogger *log.Logger
	errorLogger *log.Logger
	level       Level
}

func New(level Level) *Logger {
	debugLogger := log.New(os.Stdout, "DEBUG: [TraefikRedirector] ", 0)
	errorLogger := log.New(os.Stdout, "ERROR: [TraefikRedirector] ", 0)

	return &Logger{
		debugLogger: debugLogger,
		errorLogger: errorLogger,
		level:       level,
	}
}

func (log *Logger) Error(v ...any) {
	log.errorLogger.Println(v)
}

func (log *Logger) Debug(v ...any) {
	if log.level == Debug {
		log.debugLogger.Println(v)
	}
}
