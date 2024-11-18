package kem_x25519

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/kem"
)

func TestX25519(t *testing.T) {
	kem.TestScheme256[PrivateKey, PublicKey](t, New())
}
