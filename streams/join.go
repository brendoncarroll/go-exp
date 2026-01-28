package streams

import (
	"context"
	"fmt"

	"go.brendoncarroll.net/exp/maybe"
)

// OJoined is an outer-joined pair
type OJoined[L, R any] struct {
	Left  maybe.Maybe[L]
	Right maybe.Maybe[R]
}

// Reset sets Left and Right to Nothing
func (oj *OJoined[L, R]) Reset() {
	oj.Left.Ok = false
	oj.Right.Ok = false
}

func (oj *OJoined[L, R]) SetLeft(l L) {
	oj.Left = maybe.Just(l)
}

func (oj *OJoined[L, R]) SetRight(r R) {
	oj.Right = maybe.Just(r)
}

func (oj OJoined[L, R]) String() string {
	switch {
	case oj.Left.Ok && oj.Right.Ok:
		return fmt.Sprintf("{L: %v, R: %v}", oj.Left.X, oj.Right.X)
	case oj.Left.Ok:
		return fmt.Sprintf("{L: %v}", oj.Left.X)
	case oj.Right.Ok:
		return fmt.Sprintf("{R: %v}", oj.Right.X)
	default:
		return "{}"
	}
}

type OJoiner[L, R any] struct {
	lit Peekable[L]
	rit Peekable[R]
	cmp func(L, R) int
}

// NewOJoiner returns an Iterator that performs an outer join
// on lit and rit, which are both assumed to be sorted in increasing order.
func NewOJoiner[L, R any](lit Peekable[L], rit Peekable[R], compare func(L, R) int) *OJoiner[L, R] {
	return &OJoiner[L, R]{
		lit: lit,
		rit: rit,
		cmp: compare,
	}
}

func (j *OJoiner[A, B]) Next(ctx context.Context, dsts []OJoined[A, B]) (int, error) {
	var n int
	var leftEmpty, rightEmpty bool
	for i := range dsts {
		dsts[i].Reset()
		if err := j.lit.Peek(ctx, &dsts[i].Left.X); err != nil {
			if IsEOS(err) {
				leftEmpty = true
				break
			}
			return 0, err
		}
		if err := j.rit.Peek(ctx, &dsts[i].Right.X); err != nil {
			if IsEOS(err) {
				rightEmpty = true
				break
			}
			return 0, err
		}
		c := j.cmp(dsts[i].Left.X, dsts[i].Right.X)
		switch {
		case c < 0:
			// left < right, emit left first
			if err := Skip(ctx, j.lit, 1); err != nil {
				return 0, err
			}
			dsts[i].Left.Ok = true
		case c > 0:
			// left > right, emit right first
			if err := Skip(ctx, j.rit, 1); err != nil {
				return 0, err
			}
			dsts[i].Right.Ok = true
		default:
			// they are equal, emit both
			if err := Skip(ctx, j.lit, 1); err != nil {
				return 0, err
			}
			if err := Skip(ctx, j.rit, 1); err != nil {
				return 0, err
			}
			dsts[i].Left.Ok = true
			dsts[i].Right.Ok = true
		}
		n++
	}
	dsts2 := dsts[n:]
	// emit from right
	if leftEmpty {
		for i := range dsts2 {
			if err := NextUnit(ctx, j.rit, &dsts2[i].Right.X); err != nil {
				return n, err
			}
			dsts2[i].Right.Ok = true
			n++
		}
	}
	// emit from left
	if rightEmpty {
		for i := range dsts2 {
			if err := NextUnit(ctx, j.lit, &dsts2[i].Left.X); err != nil {
				return n, err
			}
			dsts2[i].Left.Ok = true
			n++
		}
	}
	return n, nil
}

// IJoined is the result of an inner join
type IJoined[L, R any] struct {
	Left  L
	Right R
}
