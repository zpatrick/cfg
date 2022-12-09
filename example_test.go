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
	ServerPort       cfg.Setting[int]
	ServerTimeout    cfg.Setting[time.Duration]
	DatabaseAddress  cfg.Setting[string]
	DatabaseUsername cfg.Setting[string]
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	yamlFile, err := yaml.New(path)
	if err != nil {
		return nil, err
	}

	c := &Config{
		ServerPort: cfg.Setting[int]{
			Default: cfg.Pointer(8080),
			Provider: cfg.MultiProvider[int]{
				envvar.Newf("APP_SERVER_PORT", strconv.Atoi),
				yamlFile.Int("server", "port"),
			},
		},
		ServerTimeout: cfg.Setting[time.Duration]{
			Default:   cfg.Pointer(time.Second * 30),
			Validator: cfg.Between(time.Second, time.Minute*5),
			Provider: cfg.MultiProvider[time.Duration]{
				envvar.Newf("APP_SERVER_TIMEOUT", time.ParseDuration),
				yamlFile.Duration("server", "timeout"),
			},
		},
		DatabaseAddress: cfg.Setting[string]{
			Default: cfg.Pointer("localhost:3306"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_ADDR"),
				yamlFile.String("db", "address"),
			},
		},
		DatabaseUsername: cfg.Setting[string]{
			Default:   cfg.Pointer("readonly"),
			Validator: cfg.OneOf("readonly", "readwrite"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_USERNAME"),
				yamlFile.String("db", "username"),
			},
		},
	}

	if err := cfg.Load(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func Example() {
	ctx := context.Background()
	c, err := LoadConfig(ctx, "config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mysqlConf := mysql.Config{
		Addr: c.DatabaseAddress.Val(),
		User: c.DatabaseUsername.Val(),
	}

	db, err := sql.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	svr := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.ServerPort.Val()),
		ReadTimeout:  c.ServerTimeout.Val(),
		WriteTimeout: c.ServerTimeout.Val(),
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
