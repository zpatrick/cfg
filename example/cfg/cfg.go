package cfg

import (
	"context"

	"github.com/zpatrick/cfg/example/db"
	"github.com/zpatrick/cfg/example/svr"
)

type Config struct {
	Server svr.Config
	DB     db.Config
}

func Load(ctx context.Context) (Config, error) {
	// envvar and file
}
