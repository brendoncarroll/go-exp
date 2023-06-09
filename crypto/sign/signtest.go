package sign

import (
	mrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScheme[Priv, Pub any](t *testing.T, scheme Scheme[Priv, Pub]) {
	generate := func(i int) (Pub, Priv) {
		rng := mrand.New(mrand.NewSource(int64(i)))
		pub, priv, err := scheme.Generate(rng)
		require.NoError(t, err)
		return pub, priv
	}
	t.Run("Generate", func(t *testing.T) {
		rng := mrand.New(mrand.NewSource(0))
		priv, pub, err := scheme.Generate(rng)
		require.NoError(t, err)
		require.NotNil(t, priv)
		require.NotNil(t, pub)
	})
	t.Run("DerivePublic", func(t *testing.T) {
		pub1, priv := generate(0)
		pub2 := scheme.DerivePublic(&priv)
		require.Equal(t, pub1, pub2)
	})
	t.Run("MarshalParsePublic", func(t *testing.T) {
		pub, _ := generate(0)
		data := make([]byte, scheme.PublicKeySize())
		scheme.MarshalPublic(data, &pub)
		pub2, err := scheme.ParsePublic(data)
		require.NoError(t, err)
		require.Equal(t, pub, pub2)
	})
	t.Run("MarshalParsePrivate", func(t *testing.T) {
		_, priv := generate(0)
		data := make([]byte, scheme.PrivateKeySize())
		scheme.MarshalPrivate(data, &priv)
		priv2, err := scheme.ParsePrivate(data)
		require.NoError(t, err)
		require.Equal(t, priv, priv2)
	})
	t.Run("SignVerify", func(*testing.T) {
		pub, priv := generate(0)
		sig := make([]byte, scheme.SignatureSize())
		input := []byte("hello world")
		scheme.Sign(sig, &priv, input)
		require.True(t, scheme.Verify(&pub, input, sig))

		badSig := append([]byte{}, sig...)
		badSig[0] ^= 1
		require.False(t, scheme.Verify(&pub, input, badSig))

		input2 := []byte("wrong input")
		require.False(t, scheme.Verify(&pub, input2, sig))
	})
}
