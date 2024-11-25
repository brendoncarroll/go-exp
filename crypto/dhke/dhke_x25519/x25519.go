package dhke_x25519

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"

	"go.brendoncarroll.net/exp/crypto/dhke"
)

const (
	PrivateKeySize = curve25519.ScalarSize
	PublicKeySize  = curve25519.PointSize
)

type PrivateKey = [PrivateKeySize]byte

type PublicKey = [PublicKeySize]byte

var _ dhke.Scheme[PrivateKey, PublicKey] = Scheme{}

type Scheme struct{}

func (s Scheme) Generate(rng io.Reader) (PublicKey, PrivateKey, error) {
	priv := [32]byte{}
	if _, err := io.ReadFull(rng, priv[:]); err != nil {
		return PublicKey{}, PrivateKey{}, err
	}
	pub, err := curve25519.X25519(priv[:], curve25519.Basepoint)
	return *(*[32]byte)(pub), priv, err
}

func (s Scheme) DerivePublic(priv *PrivateKey) PublicKey {
	pub, err := curve25519.X25519(priv[:], curve25519.Basepoint)
	if err != nil {
		panic(err)
	}
	return *(*[32]byte)(pub)
}

func (s Scheme) ComputeShared(dst []byte, priv *PrivateKey, pub *PublicKey) error {
	sh, err := curve25519.X25519(priv[:], pub[:])
	if err != nil {
		return err
	}
	if len(dst) != len(sh) {
		panic(fmt.Sprintf("shared is wrong length HAVE: %d WANT: %d", len(dst), len(sh)))
	}
	copy(dst[:], sh)
	return nil
}

func (s Scheme) MarshalPublic(dst []byte, x *PublicKey) {
	if len(dst) < s.PublicKeySize() {
		panic(fmt.Sprintf("len(dst) < %d", s.PublicKeySize()))
	}
	copy(dst, x[:])
}

func (s Scheme) ParsePublic(x []byte) (PublicKey, error) {
	if len(x) != 32 {
		return PublicKey{}, errors.New("wrong length for public key")
	}
	return *(*[32]byte)(x), nil
}

func (s Scheme) SharedSize() int {
	return 32
}

func (s Scheme) PublicKeySize() int {
	return 32
}

func (s Scheme) PrivateKeySize() int {
	return PrivateKeySize
}

func (s Scheme) MarshalPrivate(dst []byte, priv *PrivateKey) {
	if len(dst) < s.PrivateKeySize() {
		panic(dst)
	}
	copy(dst[:], priv[:])
}

func (s Scheme) ParsePrivate(x []byte) (PrivateKey, error) {
	if len(x) < s.PrivateKeySize() {
		return PrivateKey{}, errors.New("dhke_x25519: wrong size for private key")
	}
	var priv PrivateKey
	copy(priv[:], x)
	return priv, nil
}
