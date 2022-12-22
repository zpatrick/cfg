package cfg_test

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/cfg/providers/flags"
	"github.com/zpatrick/testx/assert"
)

func ExampleSetting_validation() {
	userName := &cfg.Setting[string]{
		Validator: cfg.OneOf("admin", "guest"),
		Provider:  cfg.StaticProvider("other"),
	}

	err := userName.Load(context.Background())
	fmt.Println(err)
	// Output: validation failed: input other not contained in [admin guest]
}

func ExampleSetting_default() {
	port := &cfg.Setting[int]{
		Default:  cfg.Pointer(8080),
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	port.Load(context.Background())
	fmt.Println(port.Val())
	// Output: 8080
}

func ExampleSetting_multiProvider() {
	addrFlag := flag.String("addr", "localhost", "")

	addr := &cfg.Setting[string]{
		Provider: cfg.MultiProvider[string]{
			envvar.New("APP_ADDR"),
			flags.NewWithDefault(addrFlag),
		},
	}

	addr.Load(context.Background())
	fmt.Println(addr.Val())
	// Output: localhost
}

// Setting.Load
// returns error if provider not defined
// returns value of Provide
// returns NoValueProvidedErorr when no default
// returns nil when NoValueProvided with default

func TestSettingLoad_returnsErrorIfProviderIsNil(t *testing.T) {
	s := cfg.Setting[int]{}
	assert.Error(t, s.Load(context.Background()))
}

func TestSettingLoad_returnsProviderValue(t *testing.T) {
	s := cfg.Setting[int]{
		Provider: cfg.ProviderFunc[int](func(ctx context.Context) (int, error) {
			return 3, nil
		}),
	}

	if err := s.Load(context.Background()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, s.Val(), 3)
}

func TestSettingLoad_returnsDefaultIfProviderReturnsNoValueProvidedErr(t *testing.T) {
	s := cfg.Setting[int]{
		Default: cfg.Pointer(3),
		Provider: cfg.ProviderFunc[int](func(ctx context.Context) (int, error) {
			return 0, cfg.NoValueProvidedError
		}),
	}

	if err := s.Load(context.Background()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, s.Val(), 3)
}
