package cfg_test

import (
	"context"
	"flag"
	"fmt"

	"github.com/zpatrick/cfg"
)

func ExampleFlag() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	if err := fs.Parse([]string{"--port", "9090"}); err != nil {
		panic(err)
	}

	port := cfg.Schema[int]{
		Name: "port",
		Providers: []cfg.Provider[int]{
			cfg.Flag(portFlag, fs, "port"),
		},
	}

	val, err := port.Load(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", val)
	// Output: port is: 9090
}

func ExampleFlagWithDefault() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	if err := fs.Parse(nil); err != nil {
		panic(err)
	}

	port := cfg.Schema[int]{
		Name: "port",
		Providers: []cfg.Provider[int]{
			cfg.FlagWithDefault(portFlag),
		},
	}

	val, err := port.Load(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", val)
	// Output: port is: 8000
}
