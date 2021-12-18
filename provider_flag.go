package cfg

import (
	"context"
	"flag"
	"strconv"
)

// A FlagLookupFunc is used to determine whether or not a given flag was set by a user.
// It's expected a FlagLookupFunc will wrap either the flag.Lookup or flag.FlagSet.Lookup methods.
// TODO: Examples:
//  portFlag := flag.Int("port", 9090, "the port to listen on")
//  cfg.IntFlag(portFlag, func() *flag.Flag { return flag.Lookup("port") }}
//  ...
//  fs := flag.NewFlagSet("", flag.ContinueOnError)
//  portFlag := fs.Int("port", 9090, "the port to listen on")
//  cfg.IntFlag(portFlag, func() *flag.Flag { return fs.Lookup("port") })
type FlagLookupFunc func() *flag.Flag

// IntFlag returns a provider from the given flag's pointer.
// If a nil lookupFunc is passed in, the flag's default value will be provided if the flag is unset.
// Otherwise, no value will be provided if the flag is unset.
func IntFlag(ptr *int, lookupFunc FlagLookupFunc) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if lookupFunc != nil && lookupFunc() == nil {
			return nil, NoValueProvidedError
		}

		s := strconv.Itoa(*ptr)
		return []byte(s), nil
	})
}
