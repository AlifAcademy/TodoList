package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"strings"
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
	StatusID    int64     `json:"status_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      int64     `json:"user_id"`
}

// Comment type
type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
}

// TagStatus type
type TagStatus struct {
	Tag    string `json:"tag"`
	Status string `json:"status"`
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
		lg.Error(err)
		return nil, err
	}

	return task, nil
}

// DeleteTaskByID method
func (s *Service) DeleteTaskByID(ctx context.Context, id int64, userID int64) (*Task, error) {
	task := &Task{}
	err := s.pool.QueryRow(ctx, `DELETE FROM tasks WHERE id=$1 and user_id=$2 RETURNING *;`, id, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return task, nil
}

// GetAllTasks method
func (s *Service) GetAllTasks(ctx context.Context, userID int64, tag string, status string, searchText string) ([]*Task, error) {
	items := make([]*Task, 0)
	var query string
	log.Println("Status:", status)
	log.Println("Tag:", tag)
	status = strings.Title(status)
	tag = strings.ToLower(tag)
	if len(status) > 0 || len(tag) > 0 {
		query = fmt.Sprintf("select t.id, t.title, t.description, t.tags, s.id status_id, t.created_at, t.updated_at, t.user_id from tasks t inner join status s on t.status_id=s.id where t.user_id=%d and (s.name='%s' or '%s' = ANY(t.tags));", userID, status, tag)
	} else if (len(searchText) > 0) {
		log.Println("Search text:", searchText)
		query = fmt.Sprintf("select * from tasks where description like '%%%s%%' and user_id=%d;", searchText, userID)
	} else {
		query = fmt.Sprintf("select * from tasks where user_id=%d", userID)
	}
	rows, err := s.pool.Query(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Task{}
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Tags, &item.StatusID, &item.CreatedAt, &item.UpdatedAt, &item.UserID)
		log.Println(item)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}

// UpdateTask method
func (s *Service) UpdateTask(ctx context.Context, item *Task, userID int64) (*Task, error) {
	task := &Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET description=$1, tags=$2 WHERE id=$3 and user_id=$4 RETURNING *;`, item.Description, item.Tags, item.ID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return task, nil
}

// MarkAsCompleted method
func (s *Service) MarkAsCompleted(ctx context.Context, taskID int64, userID int64) (*Task, error) {
	task := &Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET status_id=1 WHERE id=$1 and user_id=$2 RETURNING *;`, taskID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return task, nil
}

// MarkAsCanceled method
func (s *Service) MarkAsCanceled(ctx context.Context, taskID int64, userID int64) (*Task, error) {
	task := &Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET status_id=2 WHERE id=$1 and user_id=$2 RETURNING *;`, taskID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return task, nil
}

// AddComment method
func (s *Service) AddComment(ctx context.Context, item *Comment, userID int64) (*Comment, error) {
	comment := &Comment{}
	err := s.pool.QueryRow(ctx, `INSERT INTO comments (content, task_id, user_id) SELECT $1, $2, $3 FROM tasks WHERE id=$2 AND user_id=$3 ON CONFLICT DO NOTHING RETURNING id, content, created_at, task_id, user_id;`, item.Content, item.TaskID, userID).Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.TaskID, &comment.UserID)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return comment, nil
}

// GetUserInfo method
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (*User, error) {
	user := &User{}
	err := s.pool.QueryRow(ctx, `SELECT * FROM users WHERE id=$1;`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Hash)

	if err != nil {
		lg.Error(err)
		return nil, err
	}

	return user, nil
}

// GetStatusAndTag method
func (s *Service) GetStatusAndTag(ctx context.Context, userID int64) ([]*TagStatus, error) {
	items := make([]*TagStatus, 0)
	//query1 := fmt.Sprintf("SELECT * FROM tasks WHERE user_id=%d AND '%s' = ANY(tags)", userID, tag)

	query := fmt.Sprintf("select unnest(t.tags) as tag, s.name from tasks t inner join status s on t.status_id=s.id where t.user_id=%d;", userID)

	rows, err := s.pool.Query(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &TagStatus{}
		err := rows.Scan(&item.Tag, &item.Status)
		log.Println(item)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}