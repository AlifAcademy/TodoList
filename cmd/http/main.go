package main

import (
	"net"
	"net/http"
	"github.com/AlifAcademy/TodoList/config"
	"github.com/AlifAcademy/TodoList/internal/db/postgres"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/internal/server"
	"github.com/gorilla/mux"
	"github.com/AlifAcademy/TodoList/internal/service"
	"github.com/AlifAcademy/TodoList/internal/service/security"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/dig"
)

func main() {
	log := logger.NewFileLogger("logs.log")
	

	postgres.NewPostgresDB(config.New())

	deps := []interface{}{
		config.New,
		server.NewServer,
		mux.NewRouter,
		postgres.NewPostgresDB,
		logger.NewFileLogger,
		serverInit,
		service.NewService,
		security.NewService,
	}
	
	container := dig.New()
	for _, dep := range deps {
		err := container.Provide(dep)
		if err != nil {
			log.Error(err)
			return 
		}
	}

	err := container.Invoke(func(server *server.Server) {
		server.Init()
	})

	if err != nil {
		log.Error(err)
		return 
	}
	
	
	container.Invoke(func(server *http.Server) error {
		log.Info("Todo list starting")
		return server.ListenAndServe()
	})
	

}


func serverInit(server *server.Server, config config.Config) *http.Server{
	return &http.Server{
		Addr: net.JoinHostPort(
			config.GetString("todo.host"),
			config.GetString("todo.port"),
			),
		Handler: server,
	}
}