package server

import (
	"net/http"
	"github.com/gorilla/mux"
)
// Server type created
type Server struct {
	mux    *mux.Router
	//config *config.Config
}

// NewServer constructor
func NewServer(mux *mux.Router) *Server {
	return &Server{mux: mux}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Init Server initialization
func (s *Server) Init() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
}
