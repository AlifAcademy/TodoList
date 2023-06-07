package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/internal/middleware"
	"github.com/AlifAcademy/TodoList/internal/service"
	"github.com/AlifAcademy/TodoList/internal/service/security"
	"github.com/AlifAcademy/TodoList/pkg/types"
	"github.com/gorilla/mux"
)

var lg = logger.NewFileLogger("logs.log")

// Server type created
type Server struct {
	mux         *mux.Router
	userSvc     *service.Service
	securitySvc *security.Service
	//config *config.Config
}

const (
	// GET method
	GET = "GET"

	// POST method
	POST = "POST"

	// DELETE method
	DELETE = "DELETE"

	// UPDATE method
	UPDATE = "PUT"
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
	s.mux.Handle("/api/users", chMd(http.HandlerFunc(s.handleGetUser))).Methods(GET)

	s.mux.Handle("/api/tasks", chMd(http.HandlerFunc(s.handleNewTask))).Methods(POST)
	s.mux.Handle("/api/tasks/{id}", chMd(http.HandlerFunc(s.handleDeleteTaskByID))).Methods(DELETE)
	s.mux.Handle("/api/tasks", chMd(http.HandlerFunc(s.handleGetAllTasks))).Methods(GET)
	s.mux.Handle("/api/tasks", chMd(http.HandlerFunc(s.handleUpdateTask))).Methods(UPDATE)
	s.mux.Handle("/api/tasks/complete/{id}", chMd(http.HandlerFunc(s.handleMarkTaskAsCompeted))).Methods(UPDATE)
	s.mux.Handle("/api/tasks/cancel/{id}", chMd(http.HandlerFunc(s.handleMarkTaskAsCanceled))).Methods(UPDATE)

	s.mux.Handle("/api/comments", chMd(http.HandlerFunc(s.handleAddComment))).Methods(POST)

	s.mux.Handle("/api/tagstatus", chMd(http.HandlerFunc(s.handleGetStatusAndTag))).Methods(GET)
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

func (s *Server) handleDeleteTaskByID(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.DeleteTaskByID(request.Context(), id, userID)
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

func (s *Server) handleGetAllTasks(writer http.ResponseWriter, request *http.Request)  {
	status := request.URL.Query().Get("status")
	tag := request.URL.Query().Get("tag")
	searchText := request.URL.Query().Get("search")
	log.Print(tag)

	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.GetAllTasks(request.Context(), userID, tag, status, searchText)
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

func (s *Server) handleUpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task *service.Task
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := s.userSvc.UpdateTask(request.Context(), task, userID)

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

func (s *Server) handleMarkTaskAsCompeted(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.MarkAsCompleted(request.Context(), id, userID)
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

func (s *Server) handleMarkTaskAsCanceled(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.MarkAsCanceled(request.Context(), id, userID)
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

func (s *Server) handleAddComment(writer http.ResponseWriter, request *http.Request) {
	var comment *service.Comment
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := s.userSvc.AddComment(request.Context(), comment, userID)

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

func (s *Server) handleGetUser(writer http.ResponseWriter, request *http.Request)  {
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.GetUserInfo(request.Context(), userID)
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

func (s *Server) handleGetStatusAndTag(writer http.ResponseWriter, request *http.Request)  {
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	
	items, err := s.userSvc.GetStatusAndTag(request.Context(), userID)
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

