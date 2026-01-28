package streams

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"go.brendoncarroll.net/exp/maybe"
	"go.brendoncarroll.net/exp/slices2"
)

func TestSlice(t *testing.T) {
	ctx := context.Background()
	it := NewSlice([]int{0, 1, 2, 3, 4}, nil)

	var dst int
	for i := range 5 {
		require.NoError(t, NextUnit(ctx, it, &dst))
		require.Equal(t, i, dst)
	}
	for range 3 {
		require.ErrorIs(t, NextUnit(ctx, it, &dst), EOS())
	}
}

func TestSeq(t *testing.T) {
	type testCase = []int
	tcs := []testCase{
		{},
		{0},
		{1, 4, 10},
		{1, -1, 2, -2},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ctx := context.TODO()
			it := NewSeq(slices.Values(tc))
			defer it.Drop()
			actual, err := Collect(ctx, it, len(tc)+1)
			require.NoError(t, err)
			if len(actual) == 0 && len(tc) == 0 {
				return
			}
			require.Equal(t, tc, actual)
		})
	}
}

func TestMerge(t *testing.T) {
	type testCase struct {
		Ins [][]int
		Out []int
	}
	tcs := []testCase{
		{
			Ins: [][]int{},
			Out: nil,
		},
		{
			Ins: [][]int{
				{1, 2, 3},
			},
			Out: []int{1, 2, 3},
		},
		{
			Ins: [][]int{
				{1, 2, 3},
				{1, 2, 3},
				{1, 2, 3},
			},
			Out: []int{1, 1, 1, 2, 2, 2, 3, 3, 3},
		},
		{
			Ins: [][]int{
				{0, 2, 4, 6, 8},
				{1},
				{3, 5, 7, 9},
			},
			Out: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ctx := context.TODO()
			ins := slices2.Map(tc.Ins, func(x []int) Peekable[int] {
				return NewSlice(x, nil)
			})
			m := NewMerger(ins, cmp.Compare[int])
			actual, err := Collect(ctx, m, 100)
			require.NoError(t, err)

			require.Equal(t, tc.Out, actual)
		})
	}
}

func TestOJoiner(t *testing.T) {
	type testCase struct {
		Left  []int
		Right []int
		Out   []OJoined[int, int]
	}
	tcs := []testCase{
		{Left: nil, Right: nil,
			Out: nil,
		},
		{Left: []int{1, 10}, Right: nil,
			Out: []OJoined[int, int]{
				leftOnly(1),
				leftOnly(10),
			},
		},
		{Left: nil, Right: []int{1, 10},
			Out: []OJoined[int, int]{
				rightOnly(1),
				rightOnly(10),
			},
		},
		{Left: []int{1, 2, 3}, Right: []int{2, 3, 4},
			Out: []OJoined[int, int]{
				leftOnly(1),
				both(2),
				both(3),
				rightOnly(4),
			},
		},
		{Left: []int{1, 2, 3}, Right: []int{1, 2, 3},
			Out: []OJoined[int, int]{
				both(1),
				both(2),
				both(3),
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ctx := context.TODO()
			l := NewSlice(tc.Left, nil)
			r := NewSlice(tc.Right, nil)
			j := NewOJoiner(l, r, cmp.Compare[int])
			actual, err := Collect(ctx, j, len(tc.Left)+len(tc.Right))
			require.NoError(t, err)
			// to make equality checking easier, fully zero Nothings.
			for i := range actual {
				if !actual[i].Left.Ok {
					actual[i].Left = maybe.Nothing[int]()
				}
				if !actual[i].Right.Ok {
					actual[i].Right = maybe.Nothing[int]()
				}
			}
			require.Equal(t, tc.Out, actual)
		})
	}
}

func leftOnly[T any](x T) OJoined[T, T] {
	return OJoined[T, T]{Left: maybe.Just(x)}
}

func rightOnly[T any](x T) OJoined[T, T] {
	return OJoined[T, T]{Right: maybe.Just(x)}
}

func both[T any](x T) OJoined[T, T] {
	return OJoined[T, T]{Left: maybe.Just(x), Right: maybe.Just(x)}
}
