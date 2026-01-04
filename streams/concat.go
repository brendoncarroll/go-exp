package streams

import "context"

type concat[T any] []Iterator[T]

func (it *concat[T]) Next(ctx context.Context, dst []T) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}
	if len(*it) == 0 {
		return 0, EOS()
	}
	if err := NextUnit(ctx, (*it)[0], &dst[0]); err != nil {
		if IsEOS(err) {
			*it = (*it)[1:]
			return it.Next(ctx, dst)
		}
		return 0, err
	}
	return 0, nil
}

func Concat[T any](its ...Iterator[T]) Iterator[T] {
	c := concat[T](its)
	return &c
}
