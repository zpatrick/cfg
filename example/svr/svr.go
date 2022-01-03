package svr

import (
	"context"

	"github.com/zpatrick/cfg"
)

type Config struct {
	Port      cfg.Schema[int]
	EnableSSL cfg.Schema[bool]
}

func (c Config) Validate(ctx context.Context) error {
	return cfg.Validate(ctx, c.Port, c.EnableSSL)
}

type Server struct {
}
