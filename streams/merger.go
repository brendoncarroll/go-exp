package streams

import (
	"context"

	"go.brendoncarroll.net/exp/maybe"
)

var (
	_ Iterator[int] = &Merger[int]{}
	_ Peekable[int] = &Merger[int]{}
)

// Merger implements the merge part of the Mergesort algorithm.
type Merger[T any] struct {
	inputs []Peekable[T]
	cmp    func(a, b T) int
}

// NewMerger creates a new merging stream and returns it.
// cmp is used to determine which element should be emitted next.
func NewMerger[T any](inputs []Peekable[T], cmp func(a, b T) int) *Merger[T] {
	return &Merger[T]{
		inputs: inputs,
		cmp:    cmp,
	}
}

func (sm *Merger[T]) Next(ctx context.Context, dst []T) (int, error) {
	sr, err := sm.selectStream(ctx)
	if err != nil {
		return 0, err
	}
	if err := NextUnit(ctx, sr, &dst[0]); err != nil {
		return 0, err
	}
	return 1, nil
}

func (sm *Merger[T]) Peek(ctx context.Context, dst *T) error {
	sr, err := sm.selectStream(ctx)
	if err != nil {
		return err
	}
	return sr.Peek(ctx, dst)
}

// selectStream will never return an ended stream
func (sm *Merger[T]) selectStream(ctx context.Context) (Peekable[T], error) {
	var minTMaybe maybe.Maybe[T]
	nextIndex := len(sm.inputs)
	var ent T
	for i, sr := range sm.inputs {
		if err := sr.Peek(ctx, &ent); err != nil {
			if IsEOS(err) {
				continue
			}
			return nil, err
		}
		if !minTMaybe.Ok || sm.cmp(ent, minTMaybe.X) <= 0 {
			minTMaybe = maybe.Just(ent)
			nextIndex = i
		}
	}
	if nextIndex < len(sm.inputs) {
		return sm.inputs[nextIndex], nil
	}
	return nil, EOS()
}
