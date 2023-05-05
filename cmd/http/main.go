package main

import (
	"github.com/AlifAcademy/TodoList/internal/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.NewLogger()
	log.WithFields(log.Fields{
		"test": "another thing",
	}).Info("Just a test message 123")
}
