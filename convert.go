package cfg

import (
	"context"
	"time"
)

func Float64ToInt(f float64) (int, error) {
	return int(f), nil
}

func StringToDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func Convert[T any, A any](convert func(T) (A, error), p Provider[T]) Provider[A] {
	return ProviderFunc[A](func(ctx context.Context) (out A, err error) {
		t, err := p.Provide(ctx)
		if err != nil {
			return out, err
		}

		return convert(t)
	})
}
