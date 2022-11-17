package cfg_test

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/cfg/providers/flags"
)

func ExampleSetting_validation() {
	userName := cfg.Setting[string]{
		Name:      "UserName",
		Validator: cfg.OneOf("admin", "guest"),
		Provider:  cfg.StaticProvider("other"),
	}

	_, err := userName.Get(context.Background())
	fmt.Println(err)
	// Output: validation failed: input foo not contained in [admin guest]
}

func ExampleSetting_default() {
	port := cfg.Setting[int]{
		Name:     "port",
		Default:  cfg.Pointer(8080),
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	val, _ := port.Get(context.Background())
	fmt.Println(val)
	// Output: 8080
}

func ExampleSetting_multiProvider() {
	addrFlag := flag.String("addr", "localhost", "")

	addr := cfg.Setting[string]{
		Name: "Address",
		Provider: cfg.MultiProvider[string]{
			envvar.New("APP_ADDR"),
			flags.NewWithDefault(addrFlag),
		},
	}

	val, _ := addr.Get(context.Background())
	fmt.Println(val)
	// Output: localhost
}
