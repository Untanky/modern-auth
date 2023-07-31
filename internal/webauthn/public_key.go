package webauthn

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

type es256PublicKey struct {
	key       *ecdsa.PublicKey
	algorithm int
}

func (k *es256PublicKey) Algorithm() int {
	return -7
}

func (k *es256PublicKey) Verify(signature []byte, data []byte) bool {
	type ECDSASignature struct {
		R, S *big.Int
	}

	e := &ECDSASignature{}
	f := sha256.New
	h := f()

	h.Write(data)

	_, err := asn1.Unmarshal(signature, e)
	if err != nil {
		return false
	}

	return ecdsa.Verify(k.key, h.Sum(nil), e.R, e.S)
}
