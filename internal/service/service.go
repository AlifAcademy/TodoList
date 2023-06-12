package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"github.com/AlifAcademy/TodoList/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// ErrNotFound if an item not found
var ErrNotFound = errors.New("item not found")

// ErrInternal if some internal error occur
var ErrInternal = errors.New("internal error")

// ErrNoSuchUser if you can not find a user
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

// NewUser method
func (s *Service) NewUser(ctx context.Context, item *models.User) (*models.User, error) {
	user := &models.User{}
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
func (s *Service) NewTask(ctx context.Context, item *models.Task, userID int64) (*models.Task, error) {
	task := &models.Task{}
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
func (s *Service) DeleteTaskByID(ctx context.Context, id int64, userID int64) (*models.Task, error) {
	task := &models.Task{}
	err := s.pool.QueryRow(ctx, `DELETE FROM tasks WHERE id=$1 and user_id=$2 RETURNING *;`, id, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return task, nil
}

// GetAllTasks method
func (s *Service) GetAllTasks(ctx context.Context, userID int64, tag string, status string, searchText string) ([]*models.Task, error) {
	items := make([]*models.Task, 0)
	var query string
	log.Println("Status:", status)
	log.Println("Tag:", tag)
	status = strings.Title(status)
	tag = strings.ToLower(tag)
	if len(status) > 0 || len(tag) > 0 {
		query = fmt.Sprintf("select t.id, t.title, t.description, t.tags, s.id status_id, t.created_at, t.updated_at, t.user_id from tasks t inner join status s on t.status_id=s.id where t.user_id=%d and (s.name='%s' or '%s' = ANY(t.tags));", userID, status, tag)
	} else if len(searchText) > 0 {
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
		item := &models.Task{}
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
func (s *Service) UpdateTask(ctx context.Context, item *models.Task, userID int64) (*models.Task, error) {
	task := &models.Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET description=$1, tags=$2 WHERE id=$3 and user_id=$4 RETURNING *;`, item.Description, item.Tags, item.ID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return task, nil
}

// MarkAsCompleted method
func (s *Service) MarkAsCompleted(ctx context.Context, taskID int64, userID int64) (*models.Task, error) {
	task := &models.Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET status_id=1 WHERE id=$1 and user_id=$2 RETURNING *;`, taskID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return task, nil
}

// MarkAsCanceled method
func (s *Service) MarkAsCanceled(ctx context.Context, taskID int64, userID int64) (*models.Task, error) {
	task := &models.Task{}
	err := s.pool.QueryRow(ctx, `UPDATE tasks SET status_id=2 WHERE id=$1 and user_id=$2 RETURNING *;`, taskID, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Tags, &task.StatusID, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return task, nil
}

// AddComment method
func (s *Service) AddComment(ctx context.Context, item *models.Comment, userID int64) (*models.Comment, error) {
	comment := &models.Comment{}
	err := s.pool.QueryRow(ctx, `INSERT INTO comments (content, task_id, user_id) SELECT $1, $2, $3 FROM tasks WHERE id=$2 AND user_id=$3 ON CONFLICT DO NOTHING RETURNING id, content, created_at, task_id, user_id;`, item.Content, item.TaskID, userID).Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.TaskID, &comment.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return comment, nil
}

// GetUserInfo method
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, `SELECT * FROM users WHERE id=$1;`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Hash)

	if err != nil {
		lg.Error(err)
		return nil, ErrNoSuchUser
	}

	return user, nil
}

// GetStatusAndTag method
func (s *Service) GetStatusAndTag(ctx context.Context, userID int64) ([]*models.TagStatus, error) {
	items := make([]*models.TagStatus, 0)
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
		item := &models.TagStatus{}
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

// GetTaskByID method
func (s *Service) GetTaskByID(ctx context.Context, userID int64, taskID int64) (*models.Task, error) {
	item := &models.Task{}

	err := s.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id=$1 and user_id=$2;`, taskID, userID).Scan(&item.ID, &item.Title, &item.Description, &item.Tags, &item.StatusID, &item.CreatedAt, &item.UpdatedAt, &item.UserID)

	if err != nil {
		lg.Error(err)
		return nil, ErrNotFound
	}

	return item, nil
}

// DeleteCommentByID method
func (s *Service) DeleteCommentByID(ctx context.Context, id int64, userID int64) (*models.Comment, error) {
	comment := &models.Comment{}
	err := s.pool.QueryRow(ctx, "DELETE FROM comments WHERE id=$1 and user_id=$2 RETURNING *;", id, userID).Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.TaskID, &comment.UpdatedAt, &comment.UserID)
  
	if err != nil {
	  lg.Error(err)
	  return nil, ErrNotFound
	}
  
	return comment, nil
  }