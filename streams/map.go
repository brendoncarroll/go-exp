package streams

import (
	"context"
)

var _ Iterator[int] = &Map[float64, int]{}

type Map[X, Y any] struct {
	xs Iterator[X]
	fn func(y *Y, x X)
	x  X
}

func NewMap[X, Y any](xs Iterator[X], fn func(y *Y, x X)) *Map[X, Y] {
	return &Map[X, Y]{
		xs: xs,
		fn: fn,
	}
}

func (m Map[X, Y]) Next(ctx context.Context, dst []Y) (int, error) {
	if err := NextUnit(ctx, m.xs, &m.x); err != nil {
		return 0, err
	}
	m.fn(&dst[0], m.x)
	return 1, nil
}
