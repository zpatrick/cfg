package yaml_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal"
	"github.com/zpatrick/cfg/providers/yaml"
)

type Config struct {
	Timeout    time.Duration
	ServerPort int
	ServerAddr string
}

func Example() {
	const data = `
timeout: 5s
server:
  port: 8080
  addr: localhost
`

	path, err := internal.WriteTempFile("", data)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	yamlFile, err := yaml.New(path)
	if err != nil {
		panic(err)
	}

	var c Config
	if err := cfg.Load(context.Background(), cfg.Schemas{
		"timeout": cfg.Schema[time.Duration]{
			Dest:     &c.Timeout,
			Provider: yamlFile.Duration("timeout"),
		},
		"server.port": cfg.Schema[int]{
			Dest:     &c.ServerPort,
			Provider: yamlFile.Int("server", "port"),
		},
		"server.addr": cfg.Schema[string]{
			Dest:     &c.ServerAddr,
			Provider: yamlFile.String("server", "addr"),
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
