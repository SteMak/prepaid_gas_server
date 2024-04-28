package structs

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type Signature [65]byte

func WrapSignature(value []byte) (Signature, error) {
	var target Signature
	if len(value) != 65 {
		return target, errors.New("signature: invalid bytes length")
	}

	return *(*[65]byte)(value), nil
}

func (value Signature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
}

func (target *Signature) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))
	if err != nil {
		return errors.New("signature: unquote: " + err.Error())
	}

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) != 130 {
		return errors.New("signature: invalid length")
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return errors.New("signature: decode hex: " + err.Error())
	}
	if decoded[64] == 27 || decoded[64] == 28 {
		decoded[64] -= 27
	}

	*target, err = WrapSignature(decoded)
	return err
}

func (target *Signature) Scan(value interface{}) error {
	var err error
	*target, err = WrapSignature(value.([]byte))

	return err
}

func (sign *Signature) ToOnchain() []byte {
	data := sign[:]
	data[64] += 27
	return data
}
