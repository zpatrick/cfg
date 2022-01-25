package main

import (
	"context"
	"log"
	"math"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/example/config"
	"github.com/zpatrick/cfg/example/database"
	"github.com/zpatrick/cfg/example/server"
)

func main() {
	ctx := context.Background()
	conf, err := config.Load(ctx, config.DefaultConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.CreateDB(ctx, conf.DB)
	if err != nil {
		log.Fatal(err)
	}

	svr, err := server.CreateServer(ctx, db, conf.Server)
	if err != nil {
		log.Fatal(err)
	}

	amount := cfg.Setting[RoundedNumber]{
		Default: func() RoundedNumber { return 1 },
		Providers: []cfg.Provider[RoundedNumber]{
			LoadRoundedNumber(yamlFile, "amount"),
		},
	}

	log.Println("running service on port:", conf.Server.Port)
	log.Fatal(svr.ListenAndServe())
}

// Default your custom type.
type RoundedNumber int64

// Create a helper function for a specific provider type - we're using the cfg.YAMLFileProvider type in this example.
func ProvideRoundedNumber(provider *cfg.YAMLFileProvider, section string, keys ...string) cfg.Provider[RoundedNumber] {
	return cfg.ProviderFunc[RoundedNumber](func(ctx context.Context) (out RoundedNumber, err error) {
		// Get the raw value from the yaml provider.
		raw, err := provider.Get(section, keys...)
		if err != nil {
			return out, err
		}

		// Ensure the raw value is of the expected type.
		f, ok := raw.(float64)
		if !ok {
			return out, cfg.NewUnexpectedTypeError(float64(0), raw)
		}

		return RoundedNumber(math.Round(f)), nil
	})
}
