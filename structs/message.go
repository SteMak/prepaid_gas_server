package structs

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

type Message struct {
	Signer   []byte
	Nonce    []byte
	GasOrder []byte
	OnBehalf []byte
	Deadline uint32
	Endpoint []byte
	Gas      []byte
	Data     []byte
}

func MessageTypeHash() []byte {
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

func MessageEncode(message Message) []byte {
	pad := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	buf := []byte{}

	data_len := make([]byte, 8)
	binary.BigEndian.PutUint64(data_len, uint64(len(message.Data)))
	data_len = bytes.Join([][]byte{pad[0:24], data_len}, []byte{})

	deadline := make([]byte, 4)
	binary.BigEndian.PutUint32(deadline, message.Deadline)
	deadline = bytes.Join([][]byte{pad[0:28], deadline}, []byte{})

	buf = bytes.Join([][]byte{buf,
		pad[0:12], message.Signer,
		message.Nonce,
		message.GasOrder,
		pad[0:12], message.OnBehalf,
		deadline,
		pad[0:12], message.Endpoint,
		message.Gas,
		pad[0:31], {32}, data_len, message.Data, pad[0 : (32-len(message.Data)%32)%32],
	}, []byte{})

	return buf
}

func MessageHash(message Message) []byte {
	struct_hash := crypto.Keccak256(bytes.Join([][]byte{MessageTypeHash(), MessageEncode(message)}, []byte{}))
	domain_separator := []byte(os.Getenv("DOMAIN_SEPARATOR"))

	return crypto.Keccak256(bytes.Join([][]byte{[]byte("\x19\x01"), domain_separator, struct_hash}, []byte{}))
}
