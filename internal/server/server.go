package server

import (
	"encoding/json"
	"errors"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/internal/middleware"
	"github.com/AlifAcademy/TodoList/internal/models"
	"github.com/AlifAcademy/TodoList/internal/service"
	"github.com/AlifAcademy/TodoList/internal/service/security"
	"github.com/AlifAcademy/TodoList/pkg/types"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
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
	s.mux.Handle("/api/tasks/{id}", chMd(http.HandlerFunc(s.handleGetTaskByID))).Methods(GET)
	s.mux.Handle("/api/tasks", chMd(http.HandlerFunc(s.handleGetAllTasks))).Methods(GET)
	s.mux.Handle("/api/tasks", chMd(http.HandlerFunc(s.handleUpdateTask))).Methods(UPDATE)
	s.mux.Handle("/api/tasks/complete/{id}", chMd(http.HandlerFunc(s.handleMarkTaskAsCompeted))).Methods(UPDATE)
	s.mux.Handle("/api/tasks/cancel/{id}", chMd(http.HandlerFunc(s.handleMarkTaskAsCanceled))).Methods(UPDATE)
	s.mux.Handle("/api/comments/{id}", chMd(http.HandlerFunc(s.handleDeleteCommentByID))).Methods(DELETE)

	s.mux.Handle("/api/comments", chMd(http.HandlerFunc(s.handleAddComment))).Methods(POST)

	s.mux.Handle("/api/tagstatus", chMd(http.HandlerFunc(s.handleGetStatusAndTag))).Methods(GET)
}

func (s *Server) handleNewUser(writer http.ResponseWriter, request *http.Request) {
	var user *models.User

	err := json.NewDecoder(request.Body).Decode(&user)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}

	items, err := s.userSvc.NewUser(request.Context(), user)

	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}

	_, err = writer.Write(models.ResponseWrite("New User Successfully Created!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleNewTask(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var task *models.Task
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&task)

	if err != nil || len(task.Title) == 0 {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.StatusID = 4

	items, err := s.userSvc.NewTask(request.Context(), task, userID)

	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("New Task Successfully Created!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleDeleteTaskByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
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
		writer.Write(models.ResponseError(http.StatusNotFound, "Task Not Found").ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}

	_, err = writer.Write(models.ResponseWrite("Task Successfully Deleted!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleGetAllTasks(writer http.ResponseWriter, request *http.Request) {
	status := request.URL.Query().Get("status")
	tag := request.URL.Query().Get("tag")
	searchText := request.URL.Query().Get("search")
	log.Print(tag)

	writer.Header().Set("Content-Type", "application/json")

	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.GetAllTasks(request.Context(), userID, tag, status, searchText)
	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("Tasks successfully retrieved!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleUpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task *models.Task
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&task)

	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		lg.Error(err)
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}

	task.UpdatedAt = time.Now()

	items, err := s.userSvc.UpdateTask(request.Context(), task, userID)

	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("Task successfully updated!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleMarkTaskAsCompeted(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
		return
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.MarkAsCompleted(request.Context(), id, userID)
	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("Task marked as completed!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleMarkTaskAsCanceled(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
		return
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.MarkAsCanceled(request.Context(), id, userID)
	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("Task marked as canceled!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleAddComment(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var comment *models.Comment
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		lg.Error(err)
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}

	comment.CreatedAt = time.Now()

	items, err := s.userSvc.AddComment(request.Context(), comment, userID)

	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}

	_, err = writer.Write(models.ResponseWrite("Comment added successfully!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleGetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.GetUserInfo(request.Context(), userID)
	if errors.Is(err, service.ErrNoSuchUser) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}

	_, err = writer.Write(models.ResponseWrite("User info retrieved successfully!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleGetStatusAndTag(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.GetStatusAndTag(request.Context(), userID)
	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, http.StatusText(http.StatusNotFound)).ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}
	_, err = writer.Write(models.ResponseWrite("Tags and Statuses retrieved successfully!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleGetTaskByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		lg.Error(err)
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)

	items, err := s.userSvc.GetTaskByID(request.Context(), userID, id)
	if errors.Is(err, service.ErrNotFound) {
		writer.Write(models.ResponseError(http.StatusNotFound, "Task Not Found").ToBytes())
		return
	}
	if err != nil {
		writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
		return
	}

	_, err = writer.Write(models.ResponseWrite("Task Successfully Retrieved!", items).ToBytes())

	if err != nil {
		lg.Error(err)
		return
	}
}

func (s *Server) handleDeleteCommentByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
  
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
	  writer.Write(models.ResponseError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).ToBytes())
	  return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
	  lg.Error(err)
	}
	value := request.Context().Value(types.Key("key"))
	userID := value.(int64)
  
	items, err := s.userSvc.DeleteCommentByID(request.Context(), id, userID)
	if errors.Is(err, service.ErrNotFound) {
	  writer.Write(models.ResponseError(http.StatusNotFound, "Comment Not Found").ToBytes())
	  return
	}
	if err != nil {
	  writer.Write(models.ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).ToBytes())
	  return
	}
  
	_, err = writer.Write(models.ResponseWrite("Comment Successfully Deleted!", items).ToBytes())
  
	if err != nil {
	  lg.Error(err)
	  return
	}
  }