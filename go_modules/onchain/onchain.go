package onchain

import (
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

type Validation uint8

const (
	None Validation = iota
	StartInFuture
	NonceExhaustion
	BalanceCompliance
	OwnerCompliance
	TimelineCompliance
)

type OrderStatus uint8

const (
	None OrderStatus = iota
	Pending
	Accepted
	Active
	Inactive
	Untaken
	Closed
)

var (
	PGas *pgas.PGas

	err error
)

func Init(provider url.URL, pgas_address structs.Address, expected_separator structs.Hash) error {
	client, err := ethclient.Dial(provider.String())
	if err != nil {
		return err
	}

	address := common.BytesToAddress(pgas_address[:])
	PGas, err = pgas.NewPGas(address, client)
	if err != nil {
		return err
	}

	separator, err := PGas.DomainSeparator(nil)
	if err != nil {
		return err
	}
	if separator != expected_separator {
		return errors.New("onchain: domain separator mismatch")
	}

	return nil
}
