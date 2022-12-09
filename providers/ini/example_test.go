package ini_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/ini"
)

type Config struct {
	Timeout    cfg.Setting[time.Duration]
	ServerPort cfg.Setting[int64]
	ServerAddr cfg.Setting[string]
}

func Example() {
	const data = `
timeout = "5s"

[server]
port = 8080
addr = "localhost"
`

	path, err := cfg.WriteTempFile("", data)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	iniFile, err := ini.New(path)
	if err != nil {
		panic(err)
	}

	c := &Config{
		Timeout: cfg.Setting[time.Duration]{
			Provider: iniFile.Duration("", "timeout"),
		},
		ServerPort: cfg.Setting[int64]{
			Provider: iniFile.Int64("server", "port"),
		},
		ServerAddr: cfg.Setting[string]{
			Provider: iniFile.String("server", "addr"),
		},
	}

	if err := cfg.Load(context.Background(), c); err != nil {
		panic(err)
	}

	fmt.Printf("Timeout: %s ServerPort: %d ServerAddr: %s",
		c.Timeout.Val(),
		c.ServerPort.Val(),
		c.ServerAddr.Val(),
	)

	// Output: Timeout: 5s ServerPort: 8080 ServerAddr: localhost
}
