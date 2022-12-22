package cfg_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/cfg/providers/yaml"
)

type Config struct {
	ServerPort       int
	ServerTimeout    time.Duration
	DatabaseAddress  string
	DatabaseUsername string
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	yamlFile, err := yaml.New(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := cfg.Load(ctx, map[string]cfg.Loader{
		"server.port": cfg.Schema[int]{
			Dest:    &c.ServerPort,
			Default: cfg.Pointer(8080),
			Provider: cfg.MultiProvider[int]{
				envvar.Newf("APP_SERVER_PORT", strconv.Atoi),
				yamlFile.Int("server", "port"),
			},
		},
		"server.timeout": cfg.Schema[time.Duration]{
			Dest:      &c.ServerTimeout,
			Default:   cfg.Pointer(time.Second * 30),
			Validator: cfg.Between(time.Second, time.Minute*5),
			Provider: cfg.MultiProvider[time.Duration]{
				envvar.Newf("APP_SERVER_TIMEOUT", time.ParseDuration),
				yamlFile.Duration("server", "timeout"),
			},
		},
		"database.address": cfg.Schema[string]{
			Dest:    &c.DatabaseAddress,
			Default: cfg.Pointer("localhost:3306"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_ADDR"),
				yamlFile.String("db", "address"),
			},
		},
		"database.username": cfg.Schema[string]{
			Dest:      &c.DatabaseUsername,
			Default:   cfg.Pointer("readonly"),
			Validator: cfg.OneOf("readonly", "readwrite"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_USERNAME"),
				yamlFile.String("db", "username"),
			},
		},
	}); err != nil {
		return nil, err
	}

	return &c, nil
}

func Example() {
	ctx := context.Background()
	c, err := LoadConfig(ctx, "config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mysqlConf := mysql.Config{
		Addr: c.DatabaseAddress,
		User: c.DatabaseUsername,
	}

	db, err := sql.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	svr := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.ServerPort),
		ReadTimeout:  c.ServerTimeout,
		WriteTimeout: c.ServerTimeout,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := db.PingContext(r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
