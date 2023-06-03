package streams

import "context"

type filter[T any] struct {
	x    Iterator[T]
	pred func(T) bool
}

func NewFilter[T any](x Iterator[T], pred func(T) bool) Iterator[T] {
	return &filter[T]{}
}

func (f *filter[T]) Next(ctx context.Context, dst *T) error {
	for {
		if err := f.x.Next(ctx, dst); err != nil {
			return err
		}
		if f.pred(*dst) {
			return nil
		}
	}
}
