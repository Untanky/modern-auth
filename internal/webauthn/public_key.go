package webauthn

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"hash"
	"log"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

func decodeKey(data []byte) (PublicKey, error) {
	publicKeyData := publicKeyData{}
	err := cbor.Unmarshal(data, &publicKeyData)
	if err != nil {
		return nil, err
	}

	switch publicKeyData.KeyType {
	case 2:
		ec2PublicKey := ec2PublicKey{}
		err := cbor.Unmarshal(data, &ec2PublicKey)
		if err != nil {
			return nil, err
		}
		ec2PublicKey.publicKeyData = publicKeyData
		return &ec2PublicKey, nil
	default:
		return nil, fmt.Errorf("invalid keyType")
	}
}

type publicKeyData struct {
	KeyType int64 `cbor:"1,keyasint" json:"kty"`
	Alg     int64 `cbor:"3,keyasint" json:"alg"`
}

func (k *publicKeyData) Algorithm() int {
	return int(k.Alg)
}

func (k *publicKeyData) GetHashFunc() func() hash.Hash {
	switch k.Alg {
	case -7:
		return crypto.SHA256.New
	case -35:
		return crypto.SHA384.New
	case -36:
		return crypto.SHA512.New
	default:
		log.Printf("ERROR: invalid algorithm: %v", k.Alg)
		return nil
	}
}

type ec2PublicKey struct {
	publicKeyData
	Curve elliptic.Curve
	X     *big.Int
	Y     *big.Int
}

func (k *ec2PublicKey) UnmarshalCBOR(data []byte) error {
	type ec2Data struct {
		publicKeyData
		Curve int64  `cbor:"-1,keyasint" json:"crv"`
		X     []byte `cbor:"-2,keyasint" json:"x"`
		Y     []byte `cbor:"-3,keyasint" json:"y"`
	}
	cborData := ec2Data{}
	err := cbor.Unmarshal(data, &cborData)
	if err != nil {
		return err
	}

	switch cborData.Curve {
	case 1:
		k.Curve = elliptic.P256()
	case 2:
		k.Curve = elliptic.P384()
	case 3:
		k.Curve = elliptic.P521()
	default:
		k.Curve = elliptic.P256()
	}
	k.X = big.NewInt(0).SetBytes(cborData.X)
	k.Y = big.NewInt(0).SetBytes(cborData.Y)

	return nil
}

func (k *ec2PublicKey) Verify(signature []byte, data []byte) bool {
	key := ecdsa.PublicKey{
		Curve: k.Curve,
		X:     k.X,
		Y:     k.Y,
	}

	hash := k.GetHashFunc()()
	hash.Write(data)

	return ecdsa.VerifyASN1(&key, hash.Sum(nil), signature)
}
