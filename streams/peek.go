package streams

import (
	"context"

	"go.brendoncarroll.net/exp/maybe"
)

type Peeker[T any] struct {
	x  Iterator[T]
	cp func(dst *T, src T)

	next maybe.Maybe[T]
}

func NewPeeker[T any](x Iterator[T], cp func(dst *T, src T)) Peekable[T] {
	if p, ok := x.(Peekable[T]); ok {
		return p
	}
	if cp == nil {
		cp = func(dst *T, src T) { *dst = src }
	}
	return &Peeker[T]{
		x:  x,
		cp: cp,
	}
}

func (pi *Peeker[T]) Next(ctx context.Context, dst []T) (int, error) {
	if pi.next.Ok {
		pi.cp(&dst[0], pi.next.X)
		pi.next.Ok = false
		return 1, nil
	}
	return pi.x.Next(ctx, dst)
}

func (pi *Peeker[T]) Peek(ctx context.Context, dst *T) error {
	if !pi.next.Ok {
		if err := NextUnit(ctx, pi.x, &pi.next.X); err != nil {
			return err
		}
		pi.next.Ok = true
	}
	pi.cp(dst, pi.next.X)
	return nil
}
