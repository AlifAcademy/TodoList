package models

import "time"

// User type
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
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
	UpdatedAt   time.Time `json:"updated_at"`
}

// TagStatus type
type TagStatus struct {
	Tag    string `json:"tag"`
	Status string `json:"status"`
}
