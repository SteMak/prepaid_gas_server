package onchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/prepaidGas/prepaid-gas-server/go_modules/config"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/onchain/pgas"
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

var (
	PGas *pgas.PGas

	err error
)

func Init() error {
	client, err := ethclient.Dial(config.ProviderURL)
	if err != nil {
		return err
	}

	address := common.HexToAddress(config.PGasAddress)
	PGas, err = pgas.NewPGas(address, client)

	return err
}
