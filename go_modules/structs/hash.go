package structs

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

type Hash [32]byte

func WrapHash(value []byte) (Hash, error) {
	var target Hash
	if len(value) != 32 {
		return target, errors.New("hash: invalid bytes length")
	}

	return *(*[32]byte)(value), nil
}

func (target *Hash) Scan(value interface{}) error {
	*target, err = WrapHash(value.([]byte))
	return err
}

func (digest Hash) Sign(pkey *ecdsa.PrivateKey) (Signature, error) {
	valid_sign, err := crypto.Sign(digest[:], pkey)
	if err != nil {
		return Signature{}, errors.New("hash: sign: " + err.Error())
	}

	return WrapSignature(valid_sign)
}

func (digest Hash) Verify(sign Signature, signer Address) error {
	recovered_pubkey_bytes, err := crypto.Ecrecover(digest[:], sign[:])
	if err != nil {
		return errors.New("hash: ecrecover: " + err.Error())
	}

	recovered_pubkey, err := crypto.UnmarshalPubkey(recovered_pubkey_bytes)
	if err != nil {
		return errors.New("hash: unmarshal pkey: " + err.Error())
	}

	recovered := crypto.PubkeyToAddress(*recovered_pubkey)

	if len(recovered) != len(signer) {
		return errors.New("signature: recovered length mismatch")
	}

	for i := 0; i < len(signer); i++ {
		if signer[i] != recovered[i] {
			return errors.New("signature: recovered mismatch")
		}
	}

	return nil
}
