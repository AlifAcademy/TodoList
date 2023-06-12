package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/AlifAcademy/TodoList/config"
	"github.com/AlifAcademy/TodoList/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type configDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB new connnect to PostgreSQL Database
func NewPostgresDB(cfg config.Config) (*pgxpool.Pool, error) {

	configDb := &configDB{
		Host:     cfg.GetString("db.host"),
		Port:     cfg.GetString("db.port"),
		Username: cfg.GetString("db.username"),
		Password: cfg.GetString("db.password"),
		DBName:   cfg.GetString("db.db_name"),
		SSLMode:  "disable",
	}
	db, err := pgxpool.Connect(context.TODO(), configDb.GenerateDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	log.Println("Start seeder table create ")
	Seeder(db)

	return db, nil
}

// ------------------ Utils ------------------------ //
// GenerateDSN generate DSN string
func (c configDB) GenerateDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.DBName)
}

func Seeder(db *pgxpool.Pool) error {
	_, err := db.Exec(context.TODO(), CreateTableStatus)
	if err != nil {
	  log.Println("the Status table exists")
	}
	_, err = db.Exec(context.TODO(), CreateTableUSERS)
	if err != nil {
	  log.Println("the Users table exists")
	}
	_, err = db.Exec(context.TODO(), CreateTableTasks)
	if err != nil {
	  log.Println("the Task table exists")
	}
	_, err = db.Exec(context.TODO(), CreateTableComments)
	if err != nil {
	  log.Println("the Comments table exists")
	}

	statuses := []models.Status{
		{ID: 1, Name: "Completed", CodeName: "completed"},
		{ID: 2, Name: "Cancel", CodeName: "cancel"},
		{ID: 3, Name: "InProgress", CodeName: "in_progress"},
		{ID: 4, Name: "New", CodeName: "new"},
	}

	sqlWhere := "INSERT INTO status (id, name, code_name) VALUES"

	for i, status := range statuses {
		sqlWhere = sqlWhere + fmt.Sprintf("(%d, '%s', '%s')", status.ID, status.Name, status.CodeName)
		if i != len(statuses)-1 {
			sqlWhere += ","
		}
	}

	_, err = db.Exec(context.TODO(), sqlWhere)
	if err != nil {
		log.Println("the Statuses table exists")
	}
  
	return nil
}
