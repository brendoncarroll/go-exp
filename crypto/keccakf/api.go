package keccakf

type State1600 = [25]uint64

// KeccakF1600 applies KeccakF1600 to x inplace.
func KeccakF1600(x *State1600) {
	keccakF1600(x)
}
