package heaps

// Heap is a min heap of Ts
type Heap[T any] struct {
	lt func(a, b T) bool
	h  []T
}

// New creates a new heap
func New[T any](lt func(a, b T) bool) Heap[T] {
	return Heap[T]{lt: lt}
}

// Push adds an element to the heap
func (h *Heap[T]) Push(x T) {
	h.h = PushFunc(h.h, x, h.lt)
}

// Pop removes the minimum element from the heap and returns it.
func (h *Heap[T]) Pop() (ret T) {
	ret, h.h = PopFunc(h.h, h.lt)
	return ret
}

// Peek returns the minimum element of the heap without removing it.
func (h *Heap[T]) Peek() (ret T) {
	return Peek(h.h)
}

func (h *Heap[T]) Len() int {
	return len(h.h)
}
