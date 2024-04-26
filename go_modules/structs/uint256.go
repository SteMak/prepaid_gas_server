package structs

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
)

type Uint256 [32]byte

func WrapUint256(value []byte) (Uint256, error) {
	var target Uint256
	if len(value) != 32 {
		return target, errors.New("uint256: invalid bytes length")
	}

	return *(*[32]byte)(value), nil
}

func (value Uint256) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("0x" + hex.EncodeToString(bytes.TrimLeft(value[:], "\x00")))), nil
}

func (target *Uint256) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))
	if err != nil {
		return err
	}

	if len(hexstr) >= 2 && hexstr[0:2] == "0x" {
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

	*target, err = WrapUint256(decoded)
	return err
}

func (target *Uint256) Scan(value interface{}) error {
	len := len(value.([]byte))
	if len > 32 {
		return errors.New("uint256: invalid length")
	}

	*target, err = WrapUint256(bytes.Join([][]byte{
		pad[0 : 32-len],
		value.([]byte),
	}, []byte{}))

	return err
}

func (uint256 Uint256) IsUint32() error {
	for i := 0; i < 28; i++ {
		if uint256[i] != 0 {
			return errors.New("uint256: exceed uint32")
		}
	}

	return nil
}

func (uint256 Uint256) ToUint32() (uint32, error) {
	err = uint256.IsUint32()
	if err != nil {
		return 0, err
	}

	return uint32(uint256[31]) + uint32(uint256[30])*256 + uint32(uint256[29])*65536 + uint32(uint256[28])*16777216, nil
}

func (uint256 Uint256) ToBig() *big.Int {
	return big.NewInt(0).SetBytes(uint256[:])
}
