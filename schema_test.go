package cfg

import "testing"

func TestExample(t *testing.T) {}

// func main() {
// 	os.Setenv("APP_PORT", "9096")
//  file := cfg.NewFileProvider(cfg.FormatYAML, "config.yaml")

// 	port := Schema[int]{
// 		Name:     "port",
// 		Decode:   DecodeInt,
// 		Default:  func() int { return 9090 },
// 		Validate: Between(1, 5),
// 		Providers: []Provider{
// 			EnvVar("APP_PORT"),
// 			file.Provide("server", "port"),
// 		},
// 	}
//
// if err := Validate(port); err != nil { panic(err) }

// 	fmt.Println(port.Load())
// }
