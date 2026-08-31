package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
}

func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(os.Stdout, "[WARN]\t", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, v ...any) {
	l.infoLog.Printf(msg, v...)
}

func (l *Logger) Warn(msg string, v ...any) {
	l.warnLog.Printf(msg, v...)
}

func (l *Logger) Error(msg string, v ...any) {
	l.errorLog.Printf(msg, v...)
}
