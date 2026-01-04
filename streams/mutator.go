package streams

import (
	"context"
)

var _ Iterator[int] = &Mutator[int]{}

// Mutator edits or drops element in a stream.
// The inner stream and the Mutator contain elements of the same type.
// See Map for tranforming types.
type Mutator[T any] struct {
	x  Iterator[T]
	fn func(dst *T) bool
}

// NewMutator creates a new Mutator stream
func NewMutator[T any](x Iterator[T], fn func(dst *T) bool) *Mutator[T] {
	return &Mutator[T]{
		x:  x,
		fn: fn,
	}
}

func (fm *Mutator[T]) Next(ctx context.Context, dst []T) (int, error) {
	for {
		if err := NextUnit(ctx, fm.x, &dst[0]); err != nil {
			return 0, err
		}
		if fm.fn(&dst[0]) {
			return 0, nil
		}
	}
}
