package sig_ed25519

import (
	"testing"

	"github.com/brendoncarroll/go-exp/crypto/sign"
)

func TestEd25519(t *testing.T) {
	sign.TestScheme[PrivateKey, PublicKey](t, New())
}
