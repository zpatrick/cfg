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

	portProvider := cfg.Flag(portFlag, fs, "port")
	port, err := portProvider.Provide(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", port)
	// Output: port is: 9090
}

func ExampleFlagWithDefault() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	if err := fs.Parse(nil); err != nil {
		panic(err)
	}

	portProvider := cfg.FlagWithDefault(portFlag)
	port, err := portProvider.Provide(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", port)
	// Output: port is: 8000
}
