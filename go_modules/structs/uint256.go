package structs

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type Uint256 [32]byte

// func (value Uint256) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
// }

func (target *Uint256) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))

	if hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) > 64 {
		return errors.New("uint256: invalid length")
	}
	for len(hexstr) < 64 {
		hexstr = "0" + hexstr
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return err
	}
	if len(decoded) != 32 {
		return errors.New("uint256: invalid decode length")
	}

	*target = *(*[32]byte)(decoded)
	return nil
}

func (value Uint256) IsUint32() error {
	for i := 0; i < 28; i++ {
		if value[i] != 0 {
			return errors.New("uint256: exceed uint32")
		}
	}

	return nil
}

func (value Uint256) ToUint32() (uint32, error) {
	err = value.IsUint32()
	if err != nil {
		return 0, err
	}

	return uint32(value[31]) + uint32(value[30])*256 + uint32(value[29])*65536 + uint32(value[28])*16777216, nil
}
