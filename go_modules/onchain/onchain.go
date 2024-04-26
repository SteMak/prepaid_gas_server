package onchain

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

type Validation uint8

const (
	NoneValidation Validation = iota
	StartInFuture
	NonceExhaustion
	BalanceCompliance
	OwnerCompliance
	TimelineCompliance
)

type OrderStatus uint8

const (
	NoneStatus OrderStatus = iota
	Pending
	Accepted
	Active
	Inactive
	Untaken
	Closed
)

var (
	ClientWS   *ethclient.Client
	ClientHTTP *ethclient.Client

	PGas       *pgas.PGas
	Transactor *bind.TransactOpts

	err error
)

func Init(provider *url.URL, pgas_address structs.Address) error {
	if ClientHTTP, err = ethclient.Dial(provider.String()); err != nil {
		return errors.New("onchain: ethclient dial error: " + err.Error())
	} else if PGas, err = pgas.NewPGas(common.BytesToAddress(pgas_address[:]), ClientHTTP); err != nil {
		return errors.New("onchain: pgas instance error: " + err.Error())
	}

	return nil
}

func InitValidator(provider *url.URL, pgas_address structs.Address, expected_separator structs.Hash) error {
	if err = Init(provider, pgas_address); err != nil {
		return err
	}

	if err = ValidateSeparator(expected_separator); err != nil {
		return err
	}

	return nil
}

func InitExecutor(provider_http *url.URL, provider_ws *url.URL, pgas_address structs.Address, pkey *ecdsa.PrivateKey, chain_id uint64) error {
	if err = Init(provider_http, pgas_address); err != nil {
		return err
	}

	if ClientWS, err = ethclient.Dial(provider_ws.String()); err != nil {
		return errors.New("onchain: ethclient dial error: " + err.Error())
	}

	Transactor, err = bind.NewKeyedTransactorWithChainID(pkey, big.NewInt(0).SetUint64(chain_id))
	if err != nil {
		return err
	}

	return nil
}

func ValidateSeparator(expected_separator structs.Hash) error {
	separator, err := PGas.DomainSeparator(nil)
	if err != nil {
		return errors.New("onchain: SC query error: " + err.Error())
	}

	if separator != expected_separator {
		return errors.New("onchain: domain separator mismatch")
	}

	return nil
}

func WrapPGasMessage(message structs.Message) pgas.Message {
	return pgas.Message{
		From:  common.Address(message.From),
		Nonce: message.Nonce.ToBig(),
		Order: message.Order.ToBig(),
		Start: message.Start.ToBig(),
		To:    common.Address(message.To),
		Gas:   message.Gas.ToBig(),
		Data:  message.Data,
	}
}
