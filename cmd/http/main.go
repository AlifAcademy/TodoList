package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlifAcademy/TodoList/config"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/internal/server"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

func main() {
	
	deps := []interface{}{
		config.New,
		logger.NewLogger,
		server.NewServer,
		mux.NewRouter,
		serverInit,
	}

	container := dig.New()

	for _, dep := range deps {
		err := container.Provide(dep)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err := container.Invoke(func(server *server.Server){
		server.Init()
	})

	if err != nil {
		fmt.Println(err)
		return 
	}

	go container.Invoke(func(server *http.Server) error {
		fmt.Println("TodoList Starting ......")
		return server.ListenAndServe()
	})

	container.Invoke(func(logger logger.Logger, server *http.Server) error {
		quit := make(chan os.Signal, 1)
 	 	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
  		<-quit

  		fmt.Println("TodoList Shutting Down")
		
		return server.Shutdown(context.TODO())
	})
	
}

func serverInit(server *server.Server, config config.Config, logger logger.Logger) *http.Server{
	return &http.Server{
		Addr: net.JoinHostPort(
			config.GetString("todo.host"),
			config.GetString("todo.port"),
			),
		Handler: server,
	}
}