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
		Validator: cfg.OneOf("admin", "guest"),
		Provider:  cfg.StaticProvider("other"),
	}

	err := userName.Load(context.Background())
	fmt.Println(err)
	// Output: validation failed: input other not contained in [admin guest]
}

func ExampleSetting_default() {
	port := cfg.Setting[int]{
		Default:  cfg.Pointer(8080),
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	port.Load(context.Background())
	fmt.Println(port.Val())
	// Output: 8080
}

func ExampleSetting_multiProvider() {
	addrFlag := flag.String("addr", "localhost", "")

	addr := cfg.Setting[string]{
		Provider: cfg.MultiProvider[string]{
			envvar.New("APP_ADDR"),
			flags.NewWithDefault(addrFlag),
		},
	}

	addr.Load(context.Background())
	fmt.Println(addr.Val())
	// Output: localhost
}
