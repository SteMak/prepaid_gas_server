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
	Treasury   *pgas.PGas
	Transactor *bind.TransactOpts
)

func InitValidator(provider *url.URL, pgas_address common.Address, expected_separator structs.Hash) error {
	if err := dialProviderHTTP(provider); err != nil {
		return err
	} else if err := connectPGas(pgas_address); err != nil {
		return err
	}

	if err := validateSeparator(expected_separator); err != nil {
		return err
	}

	return nil
}

func InitExecutor(
	provider_http, provider_ws *url.URL,
	pgas_address, treasury_address common.Address,
	pkey *ecdsa.PrivateKey,
	gasfeecap, gastipcap *int64,
	chain_id uint64,
) error {
	if err := dialProviderWS(provider_ws); err != nil {
		return err
	}

	if err := dialProviderHTTP(provider_http); err != nil {
		return err
	} else if err := connectPGas(pgas_address); err != nil {
		return err
	} else if err := connectTreasury(treasury_address); err != nil {
		return err
	}

	if err := configureTransactor(pkey, gasfeecap, gastipcap, chain_id); err != nil {
		return err
	}

	return nil
}

func dialProviderHTTP(provider *url.URL) error {
	if client, err := ethclient.Dial(provider.String()); err != nil {
		return errors.New("onchain: ethclient dial error: " + err.Error())
	} else {
		ClientHTTP = client
	}

	return nil
}

func dialProviderWS(provider *url.URL) error {
	if client, err := ethclient.Dial(provider.String()); err != nil {
		return errors.New("onchain: ethclient dial error: " + err.Error())
	} else {
		ClientWS = client
	}

	return nil
}

func connectPGas(address common.Address) error {
	if instance, err := pgas.NewPGas(address, ClientHTTP); err != nil {
		return errors.New("onchain: pgas instance error: " + err.Error())
	} else {
		PGas = instance
	}

	return nil
}

func connectTreasury(address common.Address) error {
	if instance, err := pgas.NewPGas(address, ClientHTTP); err != nil {
		return errors.New("onchain: treasury instance error: " + err.Error())
	} else {
		Treasury = instance
	}

	return nil
}

func configureTransactor(pkey *ecdsa.PrivateKey, gasfeecap, gastipcap *int64, chain_id uint64) error {
	if transactor, err := bind.NewKeyedTransactorWithChainID(pkey, big.NewInt(0).SetUint64(chain_id)); err != nil {
		return errors.New("onchain: transactor: " + err.Error())
	} else {
		Transactor = transactor
	}

	if gasfeecap != nil {
		Transactor.GasFeeCap = big.NewInt(*gasfeecap)
	}
	if gastipcap != nil {
		Transactor.GasTipCap = big.NewInt(*gastipcap)
	}

	return nil
}

func validateSeparator(expected_separator structs.Hash) error {
	if separator, err := PGas.DomainSeparator(nil); err != nil {
		return errors.New("onchain: SC query error: " + err.Error())
	} else if separator != expected_separator {
		return errors.New("onchain: domain separator mismatch")
	}

	return nil
}
