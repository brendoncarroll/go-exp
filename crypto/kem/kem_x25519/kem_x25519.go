package kem_x25519

import (
	"go.brendoncarroll.net/exp/crypto/dhke/dhke_x25519"
	"go.brendoncarroll.net/exp/crypto/dhkem"
	"go.brendoncarroll.net/exp/crypto/xof/xof_sha3"
)

const (
	PrivateKeySize = dhke_x25519.PrivateKeySize
	PublicKeySize  = dhke_x25519.PublicKeySize
)

type (
	PrivateKey = dhke_x25519.PrivateKey
	PublicKey  = dhke_x25519.PublicKey
)

func New() dhkem.Scheme256[PrivateKey, PublicKey, xof_sha3.SHAKE256State] {
	return dhkem.Scheme256[PrivateKey, PublicKey, xof_sha3.SHAKE256State]{
		DH:  dhke_x25519.Scheme{},
		XOF: xof_sha3.SHAKE256{},
	}
}
