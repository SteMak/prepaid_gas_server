package structs

import (
	"errors"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaid-gas-server/go_modules/config"
)

type Hash [32]byte

func WrapHash(value []byte) (Hash, error) {
	var target Hash
	if len(value) != 32 {
		return target, errors.New("hash: invalid bytes length")
	}

	return *(*[32]byte)(value), nil
}

func (digest Hash) Sign() (Signature, error) {
	valid_sign, err := crypto.Sign(digest[:], config.ValidatorPkey)
	if err != nil {
		return Signature{}, err
	}

	return WrapSignature(valid_sign)
}

func (digest Hash) Verify(sign Signature, signer Address) error {
	recovered_pubkey_bytes, err := crypto.Ecrecover(digest[:], sign[:])
	if err != nil {
		return err
	}

	recovered_pubkey, err := crypto.UnmarshalPubkey(recovered_pubkey_bytes)
	if err != nil {
		return err
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
