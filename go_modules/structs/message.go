package structs

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

type Message struct {
	From  Address `json:"from"`
	Nonce Uint256 `json:"nonce"`
	Order Uint256 `json:"order"`
	Start Uint256 `json:"start"`
	To    Address `json:"to"`
	Gas   Uint256 `json:"gas"`
	Data  Bytes   `json:"data"`
}

func (Message) TypeHash() []byte {
	return crypto.Keccak256([]byte(
		"Message(" +
			"address from," +
			"uint256 nonce," +
			"uint256 order," +
			"uint256 start," +
			"address to," +
			"uint256 gas," +
			"bytes data" +
			")",
	))
}

func (message Message) Encode() []byte {
	data_len := make([]byte, 32)
	binary.BigEndian.PutUint64(data_len, uint64(len(message.Data)))

	return bytes.Join([][]byte{
		pad[:12], message.From[:],
		message.Nonce[:],
		message.Order[:],
		message.Start[:],
		pad[:12], message.To[:],
		message.Gas[:],
		crypto.Keccak256(message.Data),
	}, []byte{})
}

func (message Message) DigestHash(separator Hash) (Hash, error) {
	struct_hash := crypto.Keccak256(bytes.Join([][]byte{
		message.TypeHash(),
		message.Encode(),
	}, []byte{}))

	return WrapHash(crypto.Keccak256(bytes.Join([][]byte{
		[]byte("\x19\x01"),
		separator[:],
		struct_hash,
	}, []byte{})))
}
