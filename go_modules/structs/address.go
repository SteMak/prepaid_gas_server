package structs

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type Address [20]byte

// func (value Address) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
// }

func (target *Address) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))

	if hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) != 40 {
		return errors.New("address: invalid length")
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return err
	}
	if len(decoded) != 20 {
		return errors.New("address: invalid decode length")
	}

	*target = *(*[20]byte)(decoded)
	return nil
}
