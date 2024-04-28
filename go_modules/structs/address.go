package structs

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type Address [20]byte

func WrapAddress(value []byte) (Address, error) {
	var target Address
	if len(value) != 20 {
		return target, errors.New("address: invalid bytes length")
	}

	return *(*[20]byte)(value), nil
}

func (value Address) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
}

func (target *Address) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))
	if err != nil {
		return errors.New("address: unquote: " + err.Error())
	}

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) != 40 {
		return errors.New("address: invalid length")
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return errors.New("address: decode hex: " + err.Error())
	}

	*target, err = WrapAddress(decoded)
	return err
}

func (target *Address) Scan(value interface{}) error {
	var err error
	*target, err = WrapAddress(value.([]byte))

	return err
}

func (address Address) NotZero() error {
	if address == (Address{}) {
		return errors.New("address: unexpected zero address")
	}

	return nil
}
