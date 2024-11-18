package kem_sntrup

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/kem"
)

func TestSNTRUP4591761(t *testing.T) {
	kem.TestScheme256[PrivateKey4591761, PublicKey4591761](t, New4591761())
}
