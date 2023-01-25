package toml_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal/cfgtest"
	"github.com/zpatrick/cfg/providers/toml"
)

type Config struct {
	Timeout    time.Duration
	ServerPort int64
	ServerAddr string
}

func Example() {
	const data = `
timeout = "5s"

[servers]

	[servers.alpha]
	port = 8080
	addr = "localhost"
`

	path, err := cfgtest.WriteTempFile("", data)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	tomlFile, err := toml.New(path)
	if err != nil {
		panic(err)
	}

	var c Config
	if err := cfg.Load(context.Background(), cfg.Schemas{
		"timeout": cfg.Schema[time.Duration]{
			Dest:     &c.Timeout,
			Provider: tomlFile.Duration("timeout"),
		},
		"server port": cfg.Schema[int64]{
			Dest:     &c.ServerPort,
			Provider: tomlFile.Int64("servers", "alpha", "port"),
		},
		"server addr": cfg.Schema[string]{
			Dest:     &c.ServerAddr,
			Provider: tomlFile.String("servers", "alpha", "addr"),
		},
	}); err != nil {
		panic(err)
	}

	fmt.Printf("Timeout: %s ServerPort: %d ServerAddr: %s",
		c.Timeout,
		c.ServerPort,
		c.ServerAddr,
	)

	// Output: Timeout: 5s ServerPort: 8080 ServerAddr: localhost
}
