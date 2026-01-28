package streams

import (
	"context"
	"iter"
)

type Seq[T any] struct {
	next func() (T, bool)
	stop func()
}

// NewSeqErr creates an Iterator from a iter.Seq[T]
// Drop or Closed must be called on the returned Iterator
// The returned iterator never errors until EOS.
func NewSeq[T any](seq iter.Seq[T]) *Seq[T] {
	next, stop := iter.Pull(seq)
	return &Seq[T]{next: next, stop: stop}
}

func (it *Seq[T]) Next(ctx context.Context, dst []T) (int, error) {
	var n int
	for i := range dst {
		var ok bool
		dst[i], ok = it.next()
		if !ok {
			if i > 0 {
				return i, nil
			} else {
				return 0, EOS()
			}
		}
	}
	return n, nil
}

func (it *Seq[T]) Drop() {
	it.stop()
}

func (it *Seq[T]) Close() error {
	it.Drop()
	return nil
}

type SeqErr[T any] struct {
	next func() (T, error, bool)
	stop func()
}

// NewSeqErr creates an Iterator from a iter.Seq2[T, error]
// Drop or Closed must be called on the returned Iterator
func NewSeqErr[T any](x iter.Seq2[T, error]) *SeqErr[T] {
	next, stop := iter.Pull2(x)
	return &SeqErr[T]{
		next: next,
		stop: stop,
	}
}

func (it *SeqErr[T]) Next(ctx context.Context, dst []T) (int, error) {
	for i := range dst {
		y, err, ok := it.next()
		if !ok {
			if i > 0 {
				return i, nil
			} else {
				return 0, EOS()
			}
		}
		if err != nil {
			return 0, err
		}
		dst[i] = y
	}
	return len(dst), nil
}

func (it *SeqErr[T]) Drop() {
	it.stop()
}

func (it *SeqErr[T]) Close() error {
	it.Drop()
	return nil
}
