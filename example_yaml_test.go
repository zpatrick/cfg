package cfg_test

import (
	"context"
	"fmt"
	"time"

	"github.com/zpatrick/cfg"
	"gopkg.in/yaml.v3"
)

type MyConfig struct {
	Timeout    time.Duration
	ServerPort int
	ServerAddr string
}

type yamlFile struct {
	Timeout *time.Duration `yaml:"timeout"`
	Server  struct {
		Port *int    `yaml:"port"`
		Addr *string `yaml:"addr"`
	} `yaml:"server"`
}

func ExampleYAML() {
	const data = `
timeout: 5s
server:
  port: 8080
`

	var f yamlFile
	if err := yaml.Unmarshal([]byte(data), &f); err != nil {
		panic(err)
	}

	var c MyConfig
	if err := cfg.Load(context.Background(), cfg.Schemas{
		"timeout": cfg.Schema[time.Duration]{
			Dest:     &c.Timeout,
			Provider: cfg.StaticProviderAddr(f.Timeout, false),
		},
		"server.port": cfg.Schema[int]{
			Dest:     &c.ServerPort,
			Provider: cfg.StaticProviderAddr(f.Server.Port, false),
		},
		"server.addr": cfg.Schema[string]{
			Dest:     &c.ServerAddr,
			Default:  cfg.Addr("localhost"),
			Provider: cfg.StaticProviderAddr(f.Server.Addr, false),
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
