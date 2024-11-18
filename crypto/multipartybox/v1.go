package multipartybox

import (
	"go.brendoncarroll.net/exp/crypto/aead/aead_chacha20poly1305"
	"go.brendoncarroll.net/exp/crypto/kem"
	"go.brendoncarroll.net/exp/crypto/kem/kem_sntrup"
	"go.brendoncarroll.net/exp/crypto/kem/kem_x25519"
	"go.brendoncarroll.net/exp/crypto/sign/sig_ed25519"
	"go.brendoncarroll.net/exp/crypto/xof/xof_sha3"
)

type (
	KEMPrivateKeyV1  = kem.DualKey[kem_x25519.PrivateKey, kem_sntrup.PrivateKey4591761]
	KEMPublicKeyV1   = kem.DualKey[kem_x25519.PublicKey, kem_sntrup.PublicKey4591761]
	SignPrivateKeyV1 = sig_ed25519.PrivateKey
	SignPublicKeyV1  = sig_ed25519.PublicKey
	XOFStateV1       = xof_sha3.SHAKE256State

	PrivateKeyV1 = PrivateKey[KEMPrivateKeyV1, SignPrivateKeyV1]
	PublicKeyV1  = PublicKey[KEMPublicKeyV1, SignPublicKeyV1]

	SchemeV1 = Scheme[XOFStateV1, KEMPrivateKeyV1, KEMPublicKeyV1, SignPrivateKeyV1, SignPublicKeyV1]
)

// NewV1 returns the version 1 Multiparty Box encryption scheme
func NewV1() SchemeV1 {
	return SchemeV1{
		KEM: kem.Dual256[kem_x25519.PrivateKey, kem_x25519.PublicKey, kem_sntrup.PrivateKey4591761, kem_sntrup.PublicKey4591761, xof_sha3.SHAKE256State]{
			L:   kem_x25519.New(),
			R:   kem_sntrup.New4591761(),
			XOF: xof_sha3.SHAKE256{},
		},
		Sign: sig_ed25519.New(),
		AEAD: aead_chacha20poly1305.Scheme{},
		XOF:  xof_sha3.SHAKE256{},
	}
}
