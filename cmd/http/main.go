package main

import "github.com/AlifAcademy/TodoList/internal/logger"

func main() {
	log := logger.NewFileLogger("logfile.log")
	log.Info("some test")
}
