package structs

import (
	"encoding/hex"
	"strconv"
)

type Bytes []byte

func (target *Bytes) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}

	*target, err = hex.DecodeString(string(hexstr))
	return err
}
