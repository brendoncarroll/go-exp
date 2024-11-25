package multipartybox

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"go.brendoncarroll.net/exp/crypto/aead"
	"go.brendoncarroll.net/exp/crypto/kem"
	"go.brendoncarroll.net/exp/crypto/sign"
	"go.brendoncarroll.net/exp/crypto/xof"
)

type PrivateKey[KEMPriv, SigPriv any] struct {
	KEM  KEMPriv
	Sign SigPriv
}

type PublicKey[KEMPub, SigPub any] struct {
	KEM  KEMPub
	Sign SigPub
}

type Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub any] struct {
	KEM  kem.Scheme256[KEMPriv, KEMPub]
	Sign sign.Scheme[SigPriv, SigPub]
	AEAD aead.SUV256
	XOF  xof.Scheme[XOF]
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) Generate(rng io.Reader) (retPub PublicKey[KEMPub, SigPub], retPriv PrivateKey[KEMPriv, SigPriv], _ error) {
	kemPub, kemPriv, err := s.KEM.Generate(rng)
	if err != nil {
		return retPub, retPriv, err
	}
	signPub, signPriv, err := s.Sign.Generate(rng)
	if err != nil {
		return retPub, retPriv, err
	}
	retPub = PublicKey[KEMPub, SigPub]{KEM: kemPub, Sign: signPub}
	retPriv = PrivateKey[KEMPriv, SigPriv]{KEM: kemPriv, Sign: signPriv}
	return retPub, retPriv, nil
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) DerivePublic(priv PrivateKey[KEMPriv, SigPriv]) PublicKey[KEMPub, SigPub] {
	return PublicKey[KEMPub, SigPub]{
		KEM:  s.KEM.DerivePublic(&priv.KEM),
		Sign: s.Sign.DerivePublic(&priv.Sign),
	}
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) MarshalPublic(dst []byte, pub *PublicKey[KEMPub, SigPub]) {
	s.KEM.MarshalPublic(dst[:s.KEM.PublicKeySize()], &pub.KEM)
	s.Sign.MarshalPublic(dst[s.KEM.PublicKeySize():], &pub.Sign)
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) ParsePublic(x []byte) (ret PublicKey[KEMPub, SigPub], _ error) {
	kemPub, err := s.KEM.ParsePublic(x[:s.KEM.PublicKeySize()])
	if err != nil {
		return ret, err
	}
	sigPub, err := s.Sign.ParsePublic(x[s.KEM.PublicKeySize():])
	if err != nil {
		return ret, err
	}
	return PublicKey[KEMPub, SigPub]{KEM: kemPub, Sign: sigPub}, nil
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) Encrypt(out []byte, private *PrivateKey[KEMPriv, SigPriv], pubs []*KEMPub, seed *[32]byte, ptext []byte) ([]byte, error) {
	out = appendUint32(out, uint32(s.slotSize()*len(pubs)))
	slotsBegin := len(out)
	var dek, kemSeed [32]byte
	xof.DeriveKey256(s.XOF, dek[:], seed, []byte("dek"))
	xof.DeriveKey256(s.XOF, kemSeed[:], seed, []byte("kem"))
	for _, pub := range pubs {
		var err error
		out, err = s.encryptSlot(out, private, pub, &kemSeed, &dek)
		if err != nil {
			return nil, err
		}
	}
	slotsEnd := len(out)
	out = aead.AppendSealSUV256(s.AEAD, out, &dek, ptext, out[slotsBegin:slotsEnd])
	return out, nil
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) EncryptDet(out []byte, private *PrivateKey[KEMPriv, SigPriv], pubs []*KEMPub, ptext []byte) ([]byte, error) {
	var seed [32]byte
	xof.Sum(s.XOF, seed[:], ptext)
	return s.Encrypt(out, private, pubs, &seed, ptext)
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) Decrypt(out []byte, private *PrivateKey[KEMPriv, SigPriv], writers []*SigPub, ctext []byte) (int, []byte, error) {
	m, err := ParseMessage(ctext)
	if err != nil {
		return -1, nil, err
	}
	if len(m.Slots)%s.slotSize() != 0 {
		return -1, nil, fmt.Errorf("incorrect slot size")
	}
	numSlots := len(m.Slots) / s.slotSize()
	for i := 0; i < numSlots; i++ {
		begin := i * s.slotSize()
		end := (i + 1) * s.slotSize()
		sender, dek, err := s.decryptSlot(private, writers, m.Slots[begin:end])
		if err != nil {
			continue
		}
		ptex, err := aead.AppendOpenSUV256(s.AEAD, out, dek, m.Main, m.Slots)
		return sender, ptex, err
	}
	return -1, nil, errors.New("could not decrypt message")
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) encryptSlot(out []byte, private *PrivateKey[KEMPriv, SigPriv], pub *KEMPub, seed, dek *[32]byte) ([]byte, error) {
	var ss [32]byte
	kemct := make([]byte, s.KEM.CiphertextSize())
	if err := s.KEM.Encapsulate(&ss, kemct, pub, seed); err != nil {
		return nil, err
	}
	out = append(out, kemct...)

	ptext := make([]byte, s.Sign.SignatureSize()+32)
	s.Sign.Sign(ptext[:s.Sign.SignatureSize()], &private.Sign, kemct[:])
	copy(ptext[s.Sign.SignatureSize():], dek[:])
	out = aead.AppendSealSUV256(s.AEAD, out, &ss, ptext[:], kemct)
	return out, nil
}

// decryptSlot attempts to use private to recover a shared secret from the KEM ciphertext.
// if it is successful, the remaining message is interpretted as a sealed AEAD ciphertext, containing a signature and the main DEK.
func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) decryptSlot(private *PrivateKey[KEMPriv, SigPriv], pubs []*SigPub, ctext []byte) (int, *[32]byte, error) {
	kemCtext := ctext[:s.KEM.CiphertextSize()]
	aeadCtext := ctext[s.KEM.CiphertextSize():]
	var ss [32]byte
	if err := s.KEM.Decapsulate(&ss, &private.KEM, kemCtext); err != nil {
		return -1, nil, err
	}
	ptext, err := aead.AppendOpenSUV256(s.AEAD, nil, &ss, aeadCtext, kemCtext)
	if err != nil {
		return -1, nil, err
	}
	sig := ptext[:s.Sign.SignatureSize()]
	for i, pub := range pubs {
		if s.Sign.Verify(pub, kemCtext, sig) {
			dek := ptext[s.Sign.SignatureSize():]
			if len(dek) != 32 {
				return -1, nil, errors.New("DEK is wrong length")
			}
			return i, (*[32]byte)(dek), nil
		}
	}
	return -1, nil, errors.New("could not authenticate slot")
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) slotSize() int {
	const AEADKeySize = 32
	return s.KEM.CiphertextSize() + s.Sign.SignatureSize() + AEADKeySize + s.AEAD.Overhead()
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) CiphertextSize(numReaders, ptextLen int) int {
	return ptextLen + s.Overhead(numReaders)
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) PlaintextSize(ctext []byte) (int, error) {
	m, err := ParseMessage(ctext)
	if err != nil {
		return 0, err
	}
	return len(m.Main) - s.AEAD.Overhead(), nil
}

func (s *Scheme[XOF, KEMPriv, KEMPub, SigPriv, SigPub]) Overhead(numReaders int) int {
	return 4 + s.slotSize()*numReaders + s.AEAD.Overhead()
}

type Message struct {
	Slots []byte
	Main  []byte
}

func ParseMessage(x []byte) (*Message, error) {
	if len(x) > math.MaxInt32 {
		return nil, fmt.Errorf("message is too large")
	}
	if len(x) < 4 {
		return nil, fmt.Errorf("multipartybox: too short to be message")
	}
	slotsLen := binary.BigEndian.Uint32(x[:4])
	end := 4 + slotsLen
	if int(end) > len(x) {
		return nil, fmt.Errorf("varint points out of bounds")
	}
	slots := x[4:end]
	main := x[end:]
	return &Message{
		Slots: slots,
		Main:  main,
	}, nil
}

func appendUint32(out []byte, x uint32) []byte {
	buf := [4]byte{}
	binary.BigEndian.PutUint32(buf[:], x)
	return append(out, buf[:]...)
}
