package flags_test

import (
	"context"
	"flag"
	"fmt"

	"github.com/zpatrick/cfg/providers/flags"
)

func ExampleNew() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	if err := fs.Parse([]string{"--port", "9999"}); err != nil {
		panic(err)
	}

	portProvider := flags.New(fs, portFlag, "port")
	port, err := portProvider.Provide(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", port)
	// Output: port is: 9999
}

func ExampleNewWithDefault() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	if err := fs.Parse(nil); err != nil {
		panic(err)
	}

	portProvider := flags.NewWithDefault(portFlag)
	port, err := portProvider.Provide(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", port)
	// Output: port is: 8000
}
