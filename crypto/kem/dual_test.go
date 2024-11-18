package kem_test

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/kem"
	"go.brendoncarroll.net/exp/crypto/kem/kem_sntrup"
	"go.brendoncarroll.net/exp/crypto/kem/kem_x25519"
	"go.brendoncarroll.net/exp/crypto/xof/xof_sha3"
)

func TestDual(t *testing.T) {
	s := kem.Dual256[kem_x25519.PrivateKey, kem_x25519.PublicKey, kem_sntrup.PrivateKey4591761, kem_sntrup.PublicKey4591761, xof_sha3.SHAKE256State]{
		L:   kem_x25519.New(),
		R:   kem_sntrup.New4591761(),
		XOF: xof_sha3.SHAKE256{},
	}
	type Private = kem.DualKey[kem_x25519.PrivateKey, kem_sntrup.PrivateKey4591761]
	type Public = kem.DualKey[kem_x25519.PublicKey, kem_sntrup.PublicKey4591761]
	kem.TestScheme256[Private, Public](t, s)
}
