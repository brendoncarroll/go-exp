package xof

type Scheme[State any] interface {
	// New creates a new instance of the XOF
	New() State
	// Absorb appends data to the input of the XOF by modifying s
	Absorb(s *State, data []byte)
	// Expand reads out data from the XOF, and evolves s as needed to account for the read bytes.
	Expand(s *State, data []byte)
	// Reset sets s to it's initial state.
	Reset(s *State)
}

// Sum creates a new XOF state from sch, absorbs the input and expands output into dst.
func Sum[S any](sch Scheme[S], dst []byte, in []byte) {
	x := sch.New()
	sch.Absorb(&x, in)
	sch.Expand(&x, dst)
}

// SumMany creats a new XOF state from sch, absorbs the input from all of ins, in order
// and expands the output filling dst.
func SumMany[S any](sch Scheme[S], dst []byte, ins ...[]byte) {
	x := sch.New()
	for _, in := range ins {
		sch.Absorb(&x, in)
	}
	sch.Expand(&x, dst)
}

// Sum256 is a convenience function for reading 256 bits of output from an XOF.
func Sum256[S any](sch Scheme[S], in []byte) (ret [32]byte) {
	Sum(sch, ret[:], in)
	return ret
}

// Sum512 is a convenience function for reading 512 bits of output from an XOF.
func Sum512[S any](sch Scheme[S], in []byte) (ret [64]byte) {
	Sum(sch, ret[:], in)
	return ret
}

// DeriveKey256 deterministically derives a key from base and info.
// - `dst` is filled with the output.
// - `base` must have 256 bits of entropy, and should be kept secret.  If base is weak, derived keys will be also be weak.
// - `info` can be anything, but should be distinct from other infos used with this base key. info is not secret.
func DeriveKey256[S any](sch Scheme[S], dst []byte, base *[32]byte, info []byte) {
	x := sch.New()
	sch.Absorb(&x, base[:])
	sch.Absorb(&x, info)
	sch.Expand(&x, dst[:])
}

// NewRand256 seeds a random number generator using seed and returns it.
func NewRand256[S any](sch Scheme[S], seed *[32]byte) Reader[S] {
	x := sch.New()
	sch.Absorb(&x, seed[:])
	return Reader[S]{Scheme: sch, State: &x}
}

type Writer[S any] struct {
	Scheme Scheme[S]
	State  *S
}

func (w *Writer[S]) Write(p []byte) (int, error) {
	w.Scheme.Absorb(w.State, p)
	return len(p), nil
}

type Reader[S any] struct {
	Scheme Scheme[S]
	State  *S
}

func (r *Reader[S]) Read(p []byte) (int, error) {
	r.Scheme.Expand(r.State, p)
	return len(p), nil
}

type XORExpander[T any] interface {
	XOROut(x *T, dst, src []byte)
}

func XOROut[T any](sch Scheme[T], x *T, dst, src []byte) {
	if w, ok := sch.(XORExpander[T]); ok {
		w.XOROut(x, dst, src)
		return
	}
	sch.Expand(x, dst)
	for i := range src {
		dst[i] ^= src[i]
	}
}
