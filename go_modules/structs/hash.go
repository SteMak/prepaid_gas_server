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
