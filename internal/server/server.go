package server

import (
	"net/http"

	"github.com/AlifAcademy/TodoList/config"
	"github.com/gorilla/mux"
)	

type Server struct {
	mux *mux.Router
	config config.Config
}


func NewServer(mux *mux.Router) *Server {
	return &Server{mux: mux}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init(){
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
}

func (s *Server) Shutdown() {
	s.Shutdown()
}


