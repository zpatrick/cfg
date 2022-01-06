package cfg

import (
	"context"
)

func Float64ToInt(f float64) (int, error) {
	return int(f), nil
}

func Float64ToInt64(f float64) (int64, error) {
	return int64(f), nil
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
