package streams

import (
	"context"
	"time"
)

var _ Iterator[[]int] = &Batcher[int]{}

type Batcher[T any] struct {
	inner Iterator[T]
	min   int
	dur   time.Duration
}

func NewBatcher[T any](inner Iterator[T], min int, dur time.Duration) *Batcher[T] {
	return &Batcher[T]{
		inner: inner,
		min:   min,
		dur:   dur,
	}
}

func (b *Batcher[T]) Next(ctx context.Context, dst [][]T) (int, error) {
	dst2 := &dst[0]
	*dst2 = (*dst2)[:0]
	start := time.Now()
	// TODO: need context to cancel long running calls to inner.Next
	for {
		var x T
		if err := NextUnit(ctx, b.inner, &x); err != nil {
			if IsEOS(err) && len(*dst2) > 0 {
				return 1, nil
			}
			return 0, err
		}
		*dst2 = append(*dst2, x)
		if len(*dst2) >= b.min {
			break
		}
		if now := time.Now(); now.Sub(start) >= b.dur {
			break
		}
	}
	return 1, nil
}
