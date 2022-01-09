package svr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zpatrick/cfg"
	"go.uber.org/multierr"
)

type Config struct {
	Port      cfg.Schema[int]
	EnableSSL cfg.Schema[bool]
	Timeout   cfg.Schema[time.Duration]
}

func (c Config) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "port", c.Port),
		cfg.Validate(ctx, "enable_ssl", c.EnableSSL),
		cfg.Validate(ctx, "timeout", c.Timeout),
	)
}

type Server struct {
	server *http.Server
}

func CreateServer(ctx context.Context, c Config) (*Server, error) {
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", c.Port.MustLoad(ctx)),
		ReadTimeout:  c.Timeout.MustLoad(ctx),
		WriteTimeout: c.Timeout.MustLoad(ctx),
	}

	if c.EnableSSL.MustLoad(ctx) {
		// update server.TLSConfig ...
	}

	return &Server{server}, nil
}
