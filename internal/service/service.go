package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)
// ErrNotFound if an item not found
var ErrNotFound = errors.New("item not found")

// ErrInternal if some internal error occur
var ErrInternal = errors.New("internal error")

// ErrNoSuchUser if can not find a user
var ErrNoSuchUser = errors.New("no such user")

// ErrInvalidPassword if password is incorrect
var ErrInvalidPassword = errors.New("invalid password")

var lg = logger.NewFileLogger("logs.log")

// Service type
type Service struct {
	pool *pgxpool.Pool
}

// NewService constructor
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// User type
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Hash     string `json:"password_hash"`
}

// Task type
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	StatusID   int64     `json:"status_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     int64     `json:"user_id"`
}

// NewUser method
func (s *Service) NewUser(ctx context.Context, item *User) (*User, error) {
	user := &User{}
	log.Println("Users", item)
	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	err = s.pool.QueryRow(ctx, `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id, username, email, password_hash;`, item.Username, item.Email, hash).Scan(&user.ID, &user.Username, &user.Email, &user.Hash)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return user, nil
}

// NewTask method
func (s *Service) NewTask(ctx context.Context, item *Task, userID int64) (*Task, error) {
	task := &Task{}
	log.Println("Status id", item.StatusID)
	log.Println("Title", item.Title)
	err := s.pool.QueryRow(ctx, `INSERT INTO tasks (title, description, tags, status_id, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id, title, description, tags, status_id, created_at, updated_at, user_id;`, item.Title, item.Description, item.Tags, item.StatusID, item.CreatedAt, item.UpdatedAt, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		log.Println("shit")
		lg.Error(err)
		return nil, err
	}

	return task, nil
}
