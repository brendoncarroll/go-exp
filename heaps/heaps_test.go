package heaps

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPushPop(t *testing.T) {
	xs := []int{1, 8, 5, 4, 2, 9, 3, 7, 6}

	h := New(func(a, b int) bool { return a < b })
	for _, x := range xs {
		h.Push(x)
	}

	var last int
	for i := 0; h.Len() > 0; i++ {
		x := h.Pop()
		if i > 0 {
			require.Greater(t, x, last)
		}
		last = x
	}
}
