package security

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/AlifAcademy/TodoList/internal/logger"
	"log"
	"golang.org/x/crypto/bcrypt"
)

var lg = logger.NewFileLogger("logs.log")

// Service type
type Service struct {
	pool *pgxpool.Pool
}

// NewService constructor
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Auth to validation
func (s *Service) Auth(login, password string) (id int64, ok bool) {
	var userPassword string
	var userID int64
	ctx := context.Background()

	err := s.pool.QueryRow(ctx, `SELECT id, password_hash FROM users WHERE email = $1`, login).Scan(&userID, &userPassword)

	if err != nil {
		log.Print("Foo", err)
		return -1, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		lg.Error(err)
		return -1, false
	}
	return userID, true
}