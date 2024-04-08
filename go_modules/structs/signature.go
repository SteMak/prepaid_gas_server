package structs

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signature [65]byte

func WrapSignature(value []byte) (Signature, error) {
	var target Signature
	if len(value) != 65 {
		return target, errors.New("signature: invalid bytes length")
	}
	if value[64] == 27 || value[64] == 28 {
		value[64] -= 27
	}

	return *(*[65]byte)(value), nil
}

func (value Signature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
}

func (target *Signature) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) != 130 {
		return errors.New("signature: invalid length")
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return err
	}

	*target, err = WrapSignature(decoded)
	return err
}

func (target *Signature) Scan(value interface{}) error {
	*target, err = WrapSignature(value.([]byte))
	return err
}

func (sign Signature) Verify(digest Hash, address Address) error {
	recovered_pubkey_bytes, err := crypto.Ecrecover(digest[:], sign[:])
	if err != nil {
		return err
	}

	recovered_pubkey, err := crypto.UnmarshalPubkey(recovered_pubkey_bytes)
	if err != nil {
		return err
	}

	recovered := crypto.PubkeyToAddress(*recovered_pubkey)

	if len(recovered) != len(address) {
		return errors.New("signature: recovered length mismatch")
	}

	for i := 0; i < len(address); i++ {
		if address[i] != recovered[i] {
			return errors.New("signature: recovered mismatch")
		}
	}

	return nil
}
