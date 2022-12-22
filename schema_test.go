package cfg_test

import (
	"context"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestSchemaLoad_setsDestination(t *testing.T) {
	var out int

	port := cfg.Schema[int]{
		Dest:     &out,
		Provider: cfg.StaticProvider(8080),
	}

	assert.NilError(t, port.Load(context.Background()))
	assert.Equal(t, out, 8080)
}

// func ExampleSchema_validation() {
// 	userName := &cfg.Schema[string]{
// 		Validator: cfg.OneOf("admin", "guest"),
// 		Provider:  cfg.StaticProvider("other"),
// 	}

// 	err := userName.Load(context.Background())
// 	fmt.Println(err)
// 	// Output: validation failed: input other not contained in [admin guest]
// }

// func ExampleSchema_default() {
// 	port := &cfg.Schema[int]{
// 		Default:  cfg.Pointer(8080),
// 		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
// 	}

// 	port.Load(context.Background())
// 	fmt.Println(port.Val())
// 	// Output: 8080
// }

// func ExampleSchema_multiProvider() {
// 	addrFlag := flag.String("addr", "localhost", "")

// 	addr := &cfg.Schema[string]{
// 		Provider: cfg.MultiProvider[string]{
// 			envvar.New("APP_ADDR"),
// 			flags.NewWithDefault(addrFlag),
// 		},
// 	}

// 	addr.Load(context.Background())
// 	fmt.Println(addr.Val())
// 	// Output: localhost
// }

// // Schema.Load
// // returns error if provider not defined
// // returns value of Provide
// // returns NoValueProvidedErorr when no default
// // returns nil when NoValueProvided with default

// func TestSchemaLoad_returnsErrorIfProviderIsNil(t *testing.T) {
// 	s := cfg.Schema[int]{}
// 	assert.Error(t, s.Load(context.Background()))
// }

// func TestSchemaLoad_returnsProviderValue(t *testing.T) {
// 	s := cfg.Schema[int]{
// 		Provider: cfg.ProviderFunc[int](func(ctx context.Context) (int, error) {
// 			return 3, nil
// 		}),
// 	}

// 	if err := s.Load(context.Background()); err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, s.Val(), 3)
// }

// func TestSchemaLoad_returnsDefaultIfProviderReturnsNoValueProvidedErr(t *testing.T) {
// 	s := cfg.Schema[int]{
// 		Default: cfg.Pointer(3),
// 		Provider: cfg.ProviderFunc[int](func(ctx context.Context) (int, error) {
// 			return 0, cfg.NoValueProvidedError
// 		}),
// 	}

// 	if err := s.Load(context.Background()); err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, s.Val(), 3)
// }
