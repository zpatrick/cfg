package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zpatrick/cfg/example/database"
)

type Config struct {
	Port      int
	Timeout   time.Duration
	EnableTLS bool
}

type Server struct {
	*http.Server
	db *database.DB
}

func CreateServer(ctx context.Context, db *database.DB, c Config) (*Server, error) {
	s := &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", c.Port),
			ReadTimeout:  c.Timeout,
			WriteTimeout: c.Timeout,
			IdleTimeout:  c.Timeout,
		},
		db: db,
	}

	s.Server.Handler = s
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := s.db.LoadData(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
	return
}
