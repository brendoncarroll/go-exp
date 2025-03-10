package aead_xof_test

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/aead"
	"go.brendoncarroll.net/exp/crypto/aead/aead_xof"
	"go.brendoncarroll.net/exp/crypto/xof/xof_sha3"
)

func TestScheme256(t *testing.T) {
	s := aead_xof.Scheme256[xof_sha3.SHAKE256State]{XOF: xof_sha3.SHAKE256{}}
	aead.TestSUV256(t, s)
	aead.TestK256N64(t, s)
	aead.TestK256N192(t, s)
}
