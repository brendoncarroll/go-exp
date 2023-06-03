package streams

import "context"

type concat[T any] []Iterator[T]

func (it *concat[T]) Next(ctx context.Context, dst *T) error {
	if len(*it) == 0 {
		return EOS()
	}
	if err := (*it)[0].Next(ctx, dst); err != nil {
		if IsEOS(err) {
			*it = (*it)[1:]
			return it.Next(ctx, dst)
		}
		return err
	}
	return nil
}

func Concat[T any](its ...Iterator[T]) Iterator[T] {
	c := concat[T](its)
	return &c
}
