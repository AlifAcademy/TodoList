package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/pkg/types"
	"github.com/AlifAcademy/TodoList/internal/middleware"
	"github.com/AlifAcademy/TodoList/internal/service"
	"github.com/AlifAcademy/TodoList/internal/service/security"
	"github.com/gorilla/mux"
)

var lg = logger.NewFileLogger("logs.log")

// Server type created
type Server struct {
	mux    *mux.Router
	userSvc *service.Service
	securitySvc *security.Service
	//config *config.Config
}

const (
	// GET method
	GET    = "GET"

	// POST method
	POST   = "POST"

	// DELETE method
	DELETE = "DELETE"
)

// NewServer constructor
func NewServer(mux *mux.Router, userSvc *service.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, userSvc: userSvc, securitySvc: securitySvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Init Server initialization
func (s *Server) Init() {
	chMd := middleware.Basic(s.securitySvc.Auth)

	s.mux.HandleFunc("/api/users", s.handleNewUser).Methods(POST)
	s.mux.Handle("/api/task", chMd(http.HandlerFunc(s.handleNewTask))).Methods(POST)
}

func (s *Server) handleNewUser(writer http.ResponseWriter, request *http.Request) {
	var user *service.User

	err := json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := s.userSvc.NewUser(request.Context(), user)
	
	if errors.Is(err, service.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		lg.Error(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleNewTask(writer http.ResponseWriter, request *http.Request) {
	var task *service.Task
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := s.userSvc.NewTask(request.Context(), task, userID)
	
	if errors.Is(err, service.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		lg.Error(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		lg.Error(err)
		return
	}
}
