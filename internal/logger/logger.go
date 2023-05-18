package logger

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	FOLDER_LOG_NAME  = "log"
	FILE_NAME = "log.log"
)

// Logger interface
type Logger interface {
	Info(message string)
	Error(err error)
}

// FileLogger logs inside the file
type fileLogger struct {
	logrus *log.Logger
	file *os.File
}

// NewFileLogger creates a new instance of FileLogger
func NewLogger() (Logger, error) {
	_, err := os.Stat(FOLDER_LOG_NAME);
	if err != nil {
		os.Mkdir(FOLDER_LOG_NAME, 0666)	
	}

	file, err := os.OpenFile(FOLDER_LOG_NAME +"/"+ FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// defer file.Close()
	logger := log.New()
	logger.SetOutput(file)
	logger.SetFormatter(&log.JSONFormatter{})
	
	item := &fileLogger{
		logrus: logger,
		file: file,
	}

	return item, nil
}

// Print method outputs an appropriate log according to the level
func (f *fileLogger) Print(level log.Level, message string) {
	
	f.logrus.SetLevel(level)
	switch level {
	case log.InfoLevel:
		f.logrus.WithFields(log.Fields{}).Info(message)
	case log.ErrorLevel:
		f.logrus.WithFields(log.Fields{}).Error(message)
	}
}

// Info method just logs a message
func (f *fileLogger) Info(message string) {
	f.Print(log.InfoLevel, message)
}

func (f *fileLogger) Error(err error) {
	f.Print(log.ErrorLevel, err.Error())
}

func (f *fileLogger) FileClose() {
	f.file.Close()
}
