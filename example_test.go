package cfg_test

// TODO: fix
// type Config struct {
// 	ServerPort       int
// 	ServerTimeout    time.Duration
// 	DatabaseAddress  string
// 	DatabaseUsername string
// 	DatabasePassword string
// }

// func LoadConfig(ctx context.Context, path string) (*Config, error) {
// 	yamlFile, err := yaml.New(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	errs := cfg.NewErrorAggregator()
// 	c := Config{
// 		ServerPort: cfg.Setting[int]{
// 			Name:    "Server Port",
// 			Default: cfg.Pointer(8080),
// 			Provider: cfg.MultiProvider[int]{
// 				envvar.Newf("APP_SERVER_PORT", strconv.Atoi),
// 				yamlFile.Int("server", "port"),
// 			},
// 		}.MustGet(ctx, errs),
// 		ServerTimeout: cfg.Setting[time.Duration]{
// 			Name:      "Server Timeout",
// 			Default:   cfg.Pointer(time.Second * 30),
// 			Validator: cfg.Between(time.Second, time.Minute*5),
// 			Provider: cfg.MultiProvider[time.Duration]{
// 				envvar.Newf("APP_SERVER_TIMEOUT", time.ParseDuration),
// 				yamlFile.Duration("server", "timeout"),
// 			},
// 		}.MustGet(ctx, errs),
// 		DatabaseAddress: cfg.Setting[string]{
// 			Name:    "Database Address",
// 			Default: cfg.Pointer("localhost:3306"),
// 			Provider: cfg.MultiProvider[string]{
// 				envvar.New("APP_DATABASE_ADDR"),
// 				yamlFile.String("db", "address"),
// 			},
// 		}.MustGet(ctx, errs),
// 		DatabasePassword: cfg.Setting[string]{
// 			Name: "Database Address",
// 			Provider: cfg.MultiProvider[string]{
// 				envvar.New("APP_DATABASE_ADDR"),
// 				yamlFile.String("db", "address"),
// 			},
// 			Default: cfg.Pointer("localhost:3306"),
// 		}.MustGet(ctx, errs),
// 		DatabaseUsername: cfg.Setting[string]{
// 			Name:      "Database Username",
// 			Validator: cfg.OneOf("readonly_user", "readwrite_user"),
// 			Provider: cfg.MultiProvider[string]{
// 				envvar.New("APP_DATABASE_USERNAME"),
// 				yamlFile.String("db", "username"),
// 			},
// 		}.MustGet(ctx, errs),
// 	}

// 	if err := errs.Err(); err != nil {
// 		return nil, err
// 	}

// 	return &c, nil
// }

// func Example() {
// 	ctx := context.Background()
// 	conf, err := LoadConfig(ctx, "config.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("starting server on port", conf.ServerPort)
// }
