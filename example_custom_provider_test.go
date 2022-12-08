package cfg_test

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// func ExampleProvider_custom() {
// 	env := cfg.Setting[Environment]{
// 		Default:   cfg.Pointer(Development),
// 		Validator: cfg.OneOf(Development, Staging, Production),
// 		Provider: cfg.ProviderFunc[Environment](func(context.Context) (Environment, error) {
// 			appEnv := os.Getenv("APP_ENV")
// 			if appEnv == "" {
// 				return "", cfg.NoValueProvidedError
// 			}

// 			return Environment(appEnv), nil
// 		}),
// 	}

// 	// TODO: fix
// 	val, _ := env.Get(context.Background())
// 	fmt.Println(val)
// 	// Output: development
// }
