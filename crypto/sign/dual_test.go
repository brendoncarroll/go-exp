package sign_test

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/sign"
	"go.brendoncarroll.net/exp/crypto/sign/sig_ed25519"
)

func TestDual(t *testing.T) {
	type (
		Private = sig_ed25519.PrivateKey
		Public  = sig_ed25519.PublicKey
	)
	s := sign.Dual[Private, Public, Private, Public]{
		L: sig_ed25519.New(),
		R: sig_ed25519.New(),
	}
	sign.TestScheme[sign.DualKey[Private, Private], sign.DualKey[Public, Public]](t, s)
}
