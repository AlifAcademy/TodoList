package postgres

import (
	"context"
	"fmt"

	"github.com/AlifAcademy/TodoList/config"
	"github.com/jackc/pgx/v4/pgxpool"
)


type configDB struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

// NewPostgresDB new connnect to PostgreSQL Database  
func NewPostgresDB(cfg config.Config) (*pgxpool.Pool, error) {

	var configDb *configDB 
	configDb = &configDB{
		Host: cfg.GetString("db.host"),
		Port: cfg.GetString("db.port"),
		Username: cfg.GetString("db.username"),
		Password: cfg.GetString("db.password"),
		DBName: cfg.GetString("db.db_name"),
		SSLMode: "disable",
	}
	db, err := pgxpool.Connect(context.TODO(), configDb.GenerateDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	return db, nil
}



// ------------------ Utils ------------------------ //
// GenerateDSN generate DSN string
func (c configDB) GenerateDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.DBName)
}