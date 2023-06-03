package sign_test

import (
	"testing"

	"github.com/brendoncarroll/go-exp/crypto/sign"
	"github.com/brendoncarroll/go-exp/crypto/sign/sig_ed25519"
	"github.com/brendoncarroll/go-exp/crypto/xof/xof_sha3"
)

func TestPurpose(t *testing.T) {
	type (
		Private = sig_ed25519.PrivateKey
		Public  = sig_ed25519.PublicKey
	)
	s := sign.Purpose[Private, Public, xof_sha3.SHAKE256State]{
		Scheme:  sig_ed25519.New(),
		Purpose: "test",
		XOF:     xof_sha3.SHAKE256{},
	}
	sign.TestScheme[Private, Public](t, s)
}
