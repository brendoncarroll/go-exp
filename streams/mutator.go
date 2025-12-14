package streams

import (
	"context"
)

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

func (fm *Mutator[T]) Next(ctx context.Context, dst *T) error {
	for {
		if err := fm.x.Next(ctx, dst); err != nil {
			return err
		}
		if fm.fn(dst) {
			return nil
		}
	}
}
