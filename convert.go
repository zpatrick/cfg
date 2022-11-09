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

func Convert[OldType any, NewType any](convert func(OldType) (NewType, error), p Provider[OldType]) Provider[NewType] {
	return ProviderFunc[NewType](func(ctx context.Context) (out NewType, err error) {
		t, err := p.Provide(ctx)
		if err != nil {
			return out, err
		}

		return convert(t)
	})
}
