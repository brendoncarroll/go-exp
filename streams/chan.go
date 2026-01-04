package streams

import "context"

// LoadChan loads a channel from an Iterator.
// If the context is cancelled, LoadChan returns that error.
// If it.Next errors other than EOS, LoadChan returns that error.
func LoadChan[T any](ctx context.Context, it Iterator[T], out chan<- T) error {
	for {
		x, err := Next(ctx, it)
		if err != nil {
			if IsEOS(err) {
				break
			}
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- x:
		}
	}
	return nil
}

var _ Iterator[int] = make(Chan[int])

// Chan implements a stream backed by a channel
type Chan[T any] <-chan T

func (c Chan[T]) Next(ctx context.Context, dst []T) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}
	var ok bool
	dst[0], ok = <-c
	if !ok {
		return 0, EOS()
	}
	return 1, nil
}
