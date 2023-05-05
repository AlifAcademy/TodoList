package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger() {
	filename := "logfile.log"
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err)
		return
	}
	log.SetOutput(f)
}
