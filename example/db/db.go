package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/zpatrick/cfg"
	"go.uber.org/multierr"
)

type Config struct {
	Host     cfg.Schema[string]
	Port     cfg.Schema[int]
	Username cfg.Schema[string]
	Password cfg.Schema[string]
}

func (c Config) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "host", c.Host),
		cfg.Validate(ctx, "port", c.Port),
		cfg.Validate(ctx, "username", c.Username),
		cfg.Validate(ctx, "password", c.Password),
	)
}

type DB struct {
	db *sql.DB
}

func CreateDB(ctx context.Context, c Config) (*DB, error) {
	mc := mysql.Config{
		Addr:   fmt.Sprintf("%s:%d", c.Host.MustLoad(ctx), c.Port.MustLoad(ctx)),
		User:   c.Username.MustLoad(ctx),
		Passwd: c.Password.MustLoad(ctx),
	}

	db, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}
