package structs

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type Bytes []byte

func (value Bytes) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("0x" + hex.EncodeToString(value))), nil
}

func (target *Bytes) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))
	if err != nil {
		return errors.New("bytes: unquote: " + err.Error())
	}

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}

	*target, err = hex.DecodeString(string(hexstr))
	if err != nil {
		return errors.New("bytes: decode hex: " + err.Error())
	}

	return nil
}
