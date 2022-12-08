package yaml_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zpatrick/cfg"
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

	path, err := writeTempFile("", data)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	yamlFile, err := yaml.New(path)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	errs := cfg.NewErrorAggregator()

	config := Config{
		Timeout: cfg.Setting[time.Duration]{
			Name:     "timeout",
			Provider: yamlFile.Duration("timeout"),
		}.MustGet(ctx, errs),
		ServerPort: cfg.Setting[int]{
			Name:     "port",
			Provider: yamlFile.Int("server", "port"),
		}.MustGet(ctx, errs),
		ServerAddr: cfg.Setting[string]{
			Name:     "addr",
			Provider: yamlFile.String("server", "addr"),
		}.MustGet(ctx, errs),
	}

	if err := errs.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", config)
	// Output: yaml_test.Config{Timeout:5000000000, ServerPort:8080, ServerAddr:"localhost"}
}
