package structs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaid-gas-server/go_modules/config"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/onchain/pgas"
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

func (message Message) ValidateOffchain() error {
	start, err := message.Start.ToUint32()
	if err != nil {
		return err
	}

	if int64(start) <= time.Now().Unix()+int64(config.MinStartDelay) {
		return errors.New("message: message provided lately")
	}

	empty := true
	for i := 0; i < len(message.From); i++ {
		if message.From[i] != 0 {
			empty = false
			break
		}
	}
	if empty {
		return errors.New("message: message from is empty")
	}

	return nil
}

func (message Message) ValidateOnchain() error {
	nonce := big.NewInt(0).SetBytes(message.Nonce[:])
	order := big.NewInt(0).SetBytes(message.Order[:])
	start := big.NewInt(0).SetBytes(message.Start[:])
	gas := big.NewInt(0).SetBytes(message.Gas[:])

	result, err := onchain.PGas.MessageValidate(nil, pgas.Message{
		From:  common.Address(message.From),
		Nonce: nonce,
		Order: order,
		Start: start,
		To:    common.Address(message.To),
		Gas:   gas,
		Data:  (message.Data),
	})

	switch onchain.Validation(result) {
	case onchain.StartInFuture:
		return errors.New("onchain: start not in future")
	case onchain.NonceExhaustion:
		return errors.New("onchain: nonce already used")
	case onchain.BalanceCompliance:
		return errors.New("onchain: low order balance")
	case onchain.OwnerCompliance:
		return errors.New("onchain: incompliant order owner")
	case onchain.TimelineCompliance:
		return errors.New("onchain: not in order timeline")
	}

	return err
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
		pad[0:12], message.From[:],
		message.Nonce[:],
		message.Order[:],
		message.Start[:],
		pad[0:12], message.To[:],
		message.Gas[:],
		crypto.Keccak256(message.Data),
	}, []byte{})
}

func (message Message) DigestHash() (Hash, error) {
	struct_hash := crypto.Keccak256(bytes.Join([][]byte{
		message.TypeHash(),
		message.Encode(),
	}, []byte{}))

	return WrapHash(crypto.Keccak256(bytes.Join([][]byte{
		[]byte("\x19\x01"),
		config.DomainSeparator,
		struct_hash,
	}, []byte{})))
}
