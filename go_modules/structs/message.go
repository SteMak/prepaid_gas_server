package structs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/SteMak/prepaid_gas_server/go_modules/config"
)

var (
	pad = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	err error
)

type Message struct {
	Signer   Address
	Nonce    Uint256
	GasOrder Uint256
	OnBehalf Address
	Deadline Uint256
	Endpoint Address
	Gas      Uint256
	Data     []byte
}

func (message Message) ValidateEarlyLiquidation(execution_window int64) error {
	deadline, err := message.Deadline.ToUint32()
	if err != nil {
		return err
	}

	if int64(deadline) <= time.Now().Unix()+execution_window {
		return errors.New("message: early liquidation is possible")
	}

	return nil
}

func (_ Message) TypeHash() []byte {
	return crypto.Keccak256([]byte(
		"Message(" +
			"address signer," +
			"uint256 nonce," +
			"uint256 gasOrder," +
			"address onBehalf," +
			"uint256 deadline," +
			"address endpoint," +
			"uint256 gas," +
			"bytes data" +
			")",
	))
}

func (message Message) Encode() []byte {
	buf := []byte{}

	data_len := make([]byte, 32)
	binary.BigEndian.PutUint64(data_len, uint64(len(message.Data)))

	buf = bytes.Join([][]byte{buf,
		pad[0:12], message.Signer[:],
		message.Nonce[:],
		message.GasOrder[:],
		pad[0:12], message.OnBehalf[:],
		message.Deadline[:],
		pad[0:12], message.Endpoint[:],
		message.Gas[:],
		// Address of data start in terms of current ctx (the struct)
		pad[0:30], {1, 0},
		data_len,
		message.Data, pad[0 : (32-len(message.Data)%32)%32],
	}, []byte{})

	return buf
}

func (message Message) DigestHash() []byte {
	// https://ethereum.stackexchange.com/questions/113394/how-output-of-abi-encode-calculated
	struct_hash := crypto.Keccak256(bytes.Join([][]byte{
		message.TypeHash(),
		// Address of struct start in terms of current ctx (function parameters)
		pad[0:31], {64},
		message.Encode(),
	}, []byte{}))

	return crypto.Keccak256(bytes.Join([][]byte{[]byte("\x19\x01"), config.DomainSeparator, struct_hash}, []byte{}))
}

func (message Message) Sign() ([]byte, error) {
	hash := message.DigestHash()

	signature, err := crypto.Sign(hash, config.ValidatorPkey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
