package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		Key: key,
	}
}

func (pri *PrivateKey) PrivateKeyStr() string {
	return fmt.Sprintf("%x", pri.Key.D.Bytes())
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}


func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.Key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (pub *PublicKey) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", pub.Key.X.Bytes(), pub.Key.Y.Bytes())
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() Address {
	h := sha256.Sum256(k.ToSlice())

	return AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}
func (s *Signature) SignatureStr() string {
	return fmt.Sprintf("%064x%064x", s.S, s.R)
}


func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}

func String2BigIntTuple(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

func PublicKeyFromString(s string) *PublicKey {
	x, y := String2BigIntTuple(s)
	return &PublicKey{&ecdsa.PublicKey{elliptic.P256(), &x, &y}}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *PrivateKey {
	b, _ := hex.DecodeString(s[:])
	var bi big.Int
	_ = bi.SetBytes(b)
	return &PrivateKey{&ecdsa.PrivateKey{*publicKey, &bi}}
}

func SignatureFromString(s string) *Signature {
	x, y := String2BigIntTuple(s)
	return &Signature{&x, &y}
}