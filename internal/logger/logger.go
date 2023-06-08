package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// Logger interface
type Logger interface {
	Info(message string)
	Error(err error)
}

// FileLogger logs inside the file
type FileLogger struct {
	filename string
}

// NewFileLogger creates a new instance of FileLogger
func NewFileLogger(filename string) *FileLogger {
	return &FileLogger{
		filename: filename,
	}
}

// Print method outputs an appropriate log according to the level
func (f *FileLogger) Print(level log.Level, message string) {
	file, err := os.OpenFile(f.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(level)

	switch level {
	case log.InfoLevel:
		log.WithFields(log.Fields{}).Info(message)
	case log.ErrorLevel:
		log.WithFields(log.Fields{}).Error(message)
	}
}

// Info method just logs a message
func (f *FileLogger) Info(message string) {
	f.Print(log.InfoLevel, message)
}

func (f *FileLogger) Error(err error) {
	f.Print(log.ErrorLevel, err.Error())
}
