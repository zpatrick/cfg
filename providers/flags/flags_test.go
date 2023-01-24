package flags_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/flags"
	"github.com/zpatrick/testx/assert"
)

func ExampleNew() {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "the port to listen on")

	// Simulate user passing in '--port 9999'.
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

	// Simulate user passing in no arguments.
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

func TestNewReturnsNoValueProvidedErrorWhenUnset(t *testing.T) {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	portFlag := fs.Int("port", 8000, "")

	if err := fs.Parse(nil); err != nil {
		panic(err)
	}

	provider := flags.New(fs, portFlag, "port")
	_, err := provider.Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
