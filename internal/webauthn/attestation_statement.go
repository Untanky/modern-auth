package webauthn

import (
	"fmt"
)

type packedAttestationStatemment struct {
	algorithm        int
	signature        []byte
	certificateChain [][]byte
}

func (p *packedAttestationStatemment) Verify(authenticatorData AuthData, clienDataHash []byte) error {
	if len(p.certificateChain) == 0 {
		return p.validateWithoutCert(authenticatorData, clienDataHash)
	}

	return p.validateWithCert()
}

func (p *packedAttestationStatemment) validateWithoutCert(authenticatorData AuthData, clientDataHash []byte) error {
	if p.algorithm != authenticatorData.CredentialPublicKey.Algorithm() {
		return fmt.Errorf("invalid algorithm")
	}

	verificationData := append(authenticatorData.Raw, clientDataHash...)
	ok := authenticatorData.CredentialPublicKey.Verify(p.signature, verificationData)
	if !ok {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (p *packedAttestationStatemment) validateWithCert() error {
	return fmt.Errorf("not implemented")
}
