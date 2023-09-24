package slices2

import "cmp"

// Map returns a new slice with fn applied to all the elements.
func Map[A, B any, SA ~[]A](as SA, fn func(A) B) []B {
	bs := make([]B, len(as))
	for i := range as {
		bs[i] = fn(as[i])
	}
	return bs
}

// Filter modifies xs in place to remove all elements x for which fn(x) is false.
func Filter[T any, S ~[]T](xs S, fn func(T) bool) S {
	ret := xs[:0]
	for i := range xs {
		if fn(xs[i]) {
			ret = append(ret, xs[i])
		}
	}
	return ret
}

// FoldLeft implements foldl, as described here
// https://en.wikipedia.org/wiki/Fold_(higher-order_function)
func FoldLeft[X, Acc any, S ~[]X](xs S, init Acc, fn func(Acc, X) Acc) Acc {
	acc := init
	for i := range xs {
		acc = fn(acc, xs[i])
	}
	return acc
}

// FoldRight implements foldr, as described here
// https://en.wikipedia.org/wiki/Fold_(higher-order_function)
func FoldRight[X, Acc any, S ~[]X](xs S, init Acc, fn func(Acc, X) Acc) Acc {
	acc := init
	for i := len(xs) - 1; i >= 0; i-- {
		acc = fn(acc, xs[i])
	}
	return acc
}

// DedupSorted removes duplicate items according to eq
// It doesn't actually matter how the items are sorted as long as items which could be the same are adjacent.
func DedupSorted[T comparable, S ~[]T](xs S) S {
	var deleted int
	for i := range xs {
		if i > 0 && xs[i] == xs[i-1] {
			deleted++
		} else {
			xs[i-deleted] = xs[i]
		}
	}
	return xs[:len(xs)-deleted]
}

// DedupSortedFunc removes duplicate items according to eq
// It doesn't actually matter how the items are sorted as long as items which could be the same are adjacent.
// Note that the last arg `eq` is a function that returns true for equality, not for less-than.
func DedupSortedFunc[T any, S ~[]T](xs S, eq func(a, b T) bool) S {
	var deleted int
	for i := range xs {
		if i > 0 && eq(xs[i], xs[i-1]) {
			deleted++
		} else {
			xs[i-deleted] = xs[i]
		}
	}
	return xs[:len(xs)-deleted]
}

func Merge[T cmp.Ordered, S ~[]T](dst S, left, right S) int {
	return MergeFunc(dst, left, right, func(a, b T) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	})
}

func MergeFunc[T any, S ~[]T](dst S, left, right S, cmp func(a, b T) int) (n int) {
	var l, r int
	for n < len(dst) && l < len(left) && r < len(right) {
		cmp := cmp(left[l], right[r])
		if cmp < 0 {
			dst[n] = left[l]
			n++
			l++
		} else if cmp > 0 {
			dst[n] = right[r]
			n++
			r++
		} else {
			// dst[n] = right[r] would also work
			dst[n] = left[l]
			n++
			l++
			r++
		}
	}
	// only one of these loops should run
	for n < len(dst) && l < len(left) {
		dst[n] = left[l]
		n++
		l++
	}
	for n < len(dst) && r < len(right) {
		dst[n] = right[r]
		n++
		r++
	}
	return n
}
