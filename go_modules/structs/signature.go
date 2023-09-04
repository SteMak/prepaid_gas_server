package structs

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signature [65]byte

// func (value Signature) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote("0x" + hex.EncodeToString(value[:]))), nil
// }

func (target *Signature) UnmarshalJSON(value []byte) error {
	hexstr, err := strconv.Unquote(string(value))

	if hexstr[0:2] == "0x" {
		hexstr = hexstr[2:]
	}
	if len(hexstr) != 130 {
		return errors.New("address: invalid length")
	}

	decoded, err := hex.DecodeString(string(hexstr))
	if err != nil {
		return err
	}
	if len(decoded) != 65 {
		return errors.New("address: invalid decode length")
	}

	*target = *(*[65]byte)(decoded)
	return nil
}

func (sign Signature) Verify(digest []byte, address Address) error {
	recovered_pubkey_bytes, err := crypto.Ecrecover(digest, sign[:])
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
