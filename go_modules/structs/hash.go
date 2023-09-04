package structs

import "errors"

type Hash [32]byte

func WrapHash(value []byte) (Hash, error) {
	var target Hash
	if len(value) != 32 {
		return target, errors.New("hash: invalid bytes length")
	}

	return *(*[32]byte)(value), nil
}
