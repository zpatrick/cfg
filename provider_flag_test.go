package cfg_test

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/zpatrick/cfg"
)

func ExampleIntFlag() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 9999, "the port to listen on")

	// simulate the user passing in --port=8000
	if err := fs.Parse([]string{"--port", "8000"}); err != nil {
		panic(err)
	}

	port := cfg.Schema[int]{
		Name:   "port",
		Decode: cfg.DecodeInt,
		Providers: []cfg.Provider{
			cfg.IntFlag(portFlag, nil),
		},
	}

	val, err := port.Load(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", val)
	// Output: port is: 8000
}

func ExampleIntFlag_useFlagDefault() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 9999, "the port to listen on")

	// simulate the user passing in no flags
	if err := fs.Parse([]string{}); err != nil {
		panic(err)
	}

	port := cfg.Schema[int]{
		Name:   "port",
		Decode: cfg.DecodeInt,
		Providers: []cfg.Provider{
			cfg.IntFlag(portFlag, nil),
		},
	}

	val, err := port.Load(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", val)
	// Output: port is: 9999
}

func ExampleIntFlag_ignoreFlagDefault() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 0, "the port to listen on")

	// simulate the user passing in no flags
	if err := fs.Parse([]string{}); err != nil {
		panic(err)
	}

	port := cfg.Schema[int]{
		Name:   "port",
		Decode: cfg.DecodeInt,
		Providers: []cfg.Provider{
			cfg.IntFlag(portFlag, cfg.IgnoreFlagDefault("port", fs.Visit)),
		},
	}

	if _, err := port.Load(context.Background()); err != nil {
		if errors.Is(err, cfg.NoValueProvidedError) {
			fmt.Println("no value was provided")
		}
	}

	// Output: no value was provided
}
