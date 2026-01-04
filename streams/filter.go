package streams

import "context"

type filter[T any] struct {
	x    Iterator[T]
	pred func(T) bool
}

func NewFilter[T any](x Iterator[T], pred func(T) bool) Iterator[T] {
	return &filter[T]{
		x:    x,
		pred: pred,
	}
}

func (f *filter[T]) Next(ctx context.Context, dst []T) (int, error) {
	for {
		if err := NextUnit(ctx, f.x, &dst[0]); err != nil {
			return 0, err
		}
		if f.pred(dst[0]) {
			return 1, nil
		}
	}
}
