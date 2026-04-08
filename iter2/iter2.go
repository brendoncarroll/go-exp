package iter2

import "iter"

// Map returns a Seq[Y] which contains all the elements from
// xs transformed by fn.
func Map[X, Y any](xs iter.Seq[X], fn func(X) Y) iter.Seq[Y] {
	return func(yield func(Y) bool) {
		for x := range xs {
			if !yield(fn(x)) {
				return
			}
		}
	}
}

// Flatten takes a iter.Seq[iter.Seq[T]] and returns a iter.Seq[T]
func Flatten[T any](seqs iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for seq := range seqs {
			for x := range seq {
				if !yield(x) {
					return
				}
			}
		}
	}
}

// Concat concatenates seqenences of a type T.
// It is useful for prepending a blob to a seqenence blobs
// before passing the sequence to WriteStream
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for val := range seq {
				if !yield(val) {
					return
				}
			}
		}
	}
}

// Unit returns a Seq[T] which contains one element x
func Unit[T any](x T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(x)
	}
}

// Empty returns a Seq[T] with no elements
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}
