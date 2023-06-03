package slices2

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
