package streams

import "context"

type Slice[T any] struct {
	xs  []T
	pos int
	cp  func(dst *T, src T)
}

func NewSlice[T any](xs []T, cp func(*T, T)) *Slice[T] {
	if cp == nil {
		cp = func(dst *T, src T) { *dst = src }
	}
	return &Slice[T]{xs: xs, pos: 0, cp: cp}
}

func (it *Slice[T]) Next(ctx context.Context, dst []T) (int, error) {
	if it.pos >= len(it.xs) {
		return 0, EOS()
	}
	var n int
	for ; n < len(dst) && it.pos < len(it.xs); n++ {
		it.cp(&dst[n], it.xs[it.pos])
		it.pos++
	}
	return n, nil
}

// Peek implements Peeker
func (it *Slice[T]) Peek(ctx context.Context, dst *T) error {
	if it.pos >= len(it.xs) {
		return EOS()
	}
	it.cp(dst, it.xs[it.pos])
	return nil
}

// Skip implements Skipper
func (it *Slice[T]) Skip(ctx context.Context, n int) error {
	if it.pos+n > len(it.xs) {
		return EOS()
	}
	it.pos += n
	return nil
}

func (it *Slice[T]) Reset() {
	it.pos = 0
}
