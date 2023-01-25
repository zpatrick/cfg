package ini_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal/cfgtest"
	"github.com/zpatrick/cfg/providers/ini"
)

type Config struct {
	Timeout    time.Duration
	ServerPort int64
	ServerAddr string
}

func Example() {
	const data = `
timeout = "5s"

[server]
port = 8080
addr = "localhost"
`

	path, err := cfgtest.WriteTempFile("", data)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	iniFile, err := ini.New(path)
	if err != nil {
		panic(err)
	}

	var c Config
	if err := cfg.Load(context.Background(), cfg.Schemas{
		"timeout": cfg.Schema[time.Duration]{
			Dest:     &c.Timeout,
			Provider: iniFile.Duration("", "timeout"),
		},
		"server.port": cfg.Schema[int64]{
			Dest:     &c.ServerPort,
			Provider: iniFile.Int64("server", "port"),
		},
		"server.addr": cfg.Schema[string]{
			Dest:     &c.ServerAddr,
			Provider: iniFile.String("server", "addr"),
		},
	}); err != nil {
		panic(err)
	}

	fmt.Printf("Timeout: %s ServerPort: %d ServerAddr: %s", c.Timeout, c.ServerPort, c.ServerAddr)
	// Output: Timeout: 5s ServerPort: 8080 ServerAddr: localhost
}
