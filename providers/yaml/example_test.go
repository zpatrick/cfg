package yaml_test

import (
	"time"

	"github.com/zpatrick/cfg"
)

type Config struct {
	Timeout    *cfg.Setting[time.Duration]
	ServerPort *cfg.Setting[int]
	ServerAddr *cfg.Setting[string]
}

// func Example() {
// 	const data = `
// timeout: 5s
// server:
//   port: 8080
//   addr: localhost
// `

// 	path, err := cfg.WriteTempFile("", data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer os.Remove(path)

// 	yamlFile, err := yaml.New(path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	c := &Config{
// 		Timeout: &cfg.Setting[time.Duration]{
// 			Provider: yamlFile.Duration("timeout"),
// 		},
// 		ServerPort: &cfg.Setting[int]{
// 			Provider: yamlFile.Int("server", "port"),
// 		},
// 		ServerAddr: &cfg.Setting[string]{
// 			Provider: yamlFile.String("server", "addr"),
// 		},
// 	}

// 	if err := cfg.Load(context.Background(), c); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Timeout: %s ServerPort: %d ServerAddr: %s",
// 		c.Timeout.Val(),
// 		c.ServerPort.Val(),
// 		c.ServerAddr.Val(),
// 	)

// 	// Output: Timeout: 5s ServerPort: 8080 ServerAddr: localhost
// }
