package sig_ed25519

import (
	"testing"

	"go.brendoncarroll.net/exp/crypto/sign"
)

func TestEd25519(t *testing.T) {
	sign.TestScheme[PrivateKey, PublicKey](t, New())
}
