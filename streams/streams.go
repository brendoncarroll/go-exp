package streams

import (
	"context"
	"errors"

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
	// Next advances the iterator and reads the next element into dst
	Next(ctx context.Context, dst *T) error
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
	err := it.Next(ctx, &dst)
	return dst, err
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
		if err := it.Next(ctx, &dst); err != nil {
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
	for i := range buf {
		if err := it.Next(ctx, &buf[i]); err != nil {
			return i, err
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
