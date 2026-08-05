package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, v ...interface{}) {
	l.infoLog.Printf(msg, v...)
}

func (l *Logger) Error(msg string, v ...interface{}) {
	l.errorLog.Printf(msg, v...)
}
