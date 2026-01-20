package streams

import (
	"context"
	"errors"
	"fmt"

	"go.brendoncarroll.net/exp/maybe"
)

// EndOfStream is returned by Next and Seek to indicate that the stream has no more elements.
type EndOfStream struct{}

func (EndOfStream) Error() string {
	return "end of stream"
}

// EOS returns a new EndOfStream error
func EOS() error {
	return EndOfStream{}
}

// IsEOS returns true if the error is an EndOfStream error
func IsEOS(err error) bool {
	return errors.Is(err, EndOfStream{})
}

type Iterator[T any] interface {
	// Next advances the iterator and reads the next elements into dst.
	// If len(dst) > 0, then Next must either return an error, or return n > 0.
	// That means Next must block until at least 1 element is available, if len(dst) > 0.
	// If the end of the stream has been reached, then Next returns EOS.
	// Once Next has returned EOS, it must always return (0, EOS) forever after.
	// If err != nil, then n is meaningless, and *should be* 0.
	//
	// Callers should not pass a buffer of length zero, the behavior
	// is not specified at the interface level.
	// If len(dst) == 0:
	//   - Next may return (0, nil), although this has the potential for an infinite loop.
	//   - Next may also panic if len(dst) == 0, this can be annoying, but makes misuse obvious.
	Next(ctx context.Context, dst []T) (int, error)
}

// Peekable is an Iterator which also has the Peek method
type Peekable[T any] interface {
	Iterator[T]

	// Peek shows the next element of the Iterator without changing the state of the Iterator
	Peek(ctx context.Context, dst *T) error
}

// Seeker contains the Seek method
type Seeker[T any] interface {
	// Seek ensures that all future elements of the iterator will be >= gteq
	Seek(ctx context.Context, gteq T) error
}

// Reader contains the Read method
type Reader[T any] interface {
	Read(ctx context.Context, dst []T) (int, error)
}

// Next returns a new T instead of writing it to a pointer destination.
// It calls it.Next on the Iterator
func Next[T any](ctx context.Context, it Iterator[T]) (T, error) {
	var dst T
	err := NextUnit(ctx, it, &dst)
	return dst, err
}

// NextUnit reads a single element from an Iterator
func NextUnit[T any](ctx context.Context, it Iterator[T], dst *T) error {
	var xs [1]T
	n, err := it.Next(ctx, xs[:])
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("streams: incorrect Iterator, returned n<1")
	}
	*dst = xs[0]
	return nil
}

// Peek return a new T instead of writing it to a pointer destination.
// It calls it.Peek on the Iterator.
func Peek[T any](ctx context.Context, it Peekable[T]) (T, error) {
	var dst T
	err := it.Peek(ctx, &dst)
	return dst, err
}

// ForEach calls fn for each element of it.
// fn must not retain dst, between calls.
func ForEach[T any](ctx context.Context, it Iterator[T], fn func(T) error) error {
	var dst T
	for {
		if err := NextUnit(ctx, it, &dst); err != nil {
			if IsEOS(err) {
				break
			}
			return err
		}
		if err := fn(dst); err != nil {
			return err
		}
	}
	var zero T
	dst = zero
	return nil
}

// ReadFull copies elements from the iterator into buf.
// ReadFull returns EOS when the iterator is empty.
func ReadFull[T any](ctx context.Context, it Iterator[T], buf []T) (int, error) {
	var n int
	for n < len(buf) {
		if n2, err := it.Next(ctx, buf[n:]); err != nil {
			return 0, err
		} else if n < 1 {
			return 0, fmt.Errorf("streams: incorrect iterator")
		} else {
			n += n2
		}
	}
	return len(buf), nil
}

// Collect is used to collect all of the items from an Iterator.
// If more than max elements are emitted, then Collect will return an error.
func Collect[T any](ctx context.Context, it Iterator[T], max int) (ret []T, _ error) {
	for {
		x, err := Next(ctx, it)
		if err != nil {
			if IsEOS(err) {
				break
			}
			return ret, err
		}
		if len(ret) < max {
			ret = append(ret, x)
		} else {
			return ret, errors.New("streams: too many elements to collect")
		}
	}
	return ret, nil
}

// First returns (Just, nil) when the stream produces the next element or (Nothing, nil) if the stream is over.
// First returns (Nothing, err) if it encouters any error.
func First[T any](ctx context.Context, it Iterator[T]) (maybe.Maybe[T], error) {
	x, err := Next(ctx, it)
	if err != nil {
		if IsEOS(err) {
			err = nil
		}
		return maybe.Nothing[T](), err
	}
	return maybe.Just(x), nil
}

// Last returns the last element that the stream produces.
func Last[T any](ctx context.Context, it Iterator[T]) (last maybe.Maybe[T], _ error) {
	for {
		x, err := Next(ctx, it)
		if err != nil {
			if IsEOS(err) {
				break
			}
			return maybe.Nothing[T](), err
		}
		last.X = x
		last.Ok = true
	}
	return last, nil
}
