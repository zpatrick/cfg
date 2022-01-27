package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type DB struct {
	db *sql.DB
}

func CreateDB(ctx context.Context, c Config) (*DB, error) {
	mc := mysql.Config{
		Addr:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:   c.Username,
		Passwd: c.Password,
	}

	db, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) LoadData(ctx context.Context) ([]byte, error) {
	s := fmt.Sprintf("The current time is: %s", time.Now())
	return []byte(s), nil
}
