package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"

	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

var (
	PostgresUser     string
	PostgresPassword string

	ProviderHTTP *url.URL
	ProviderWS   *url.URL
	ChainID      uint64

	DomainSeparator structs.Hash
	PGasAddress     structs.Address
	TreasuryAddress structs.Address

	ValidatorPort uint64
	ValidatorPkey *ecdsa.PrivateKey
	MinStartDelay uint64

	ExecutorPkey     *ecdsa.PrivateKey
	ExecutorAddress  common.Address
	PrevalidateDelay uint64

	err error
)

func Init() error {
	if err = godotenv.Load(); err != nil {
		return errors.New("config: environment load error: " + err.Error())
	}

	return nil
}

func InitValidator() error {
	if err = Init(); err != nil {
		return err
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	if address, err := hex.DecodeString(os.Getenv("PGAS_ADDRESS")); err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	} else if err = PGasAddress.Scan(address); err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	}

	if ProviderHTTP, err = url.Parse(os.Getenv("PROVIDER_HTTP")); err != nil {
		return errors.New("config: provider load error: " + err.Error())
	}

	if hash, err := hex.DecodeString(os.Getenv("DOMAIN_SEPARATOR")); err != nil {
		return errors.New("config: domain separator load error: " + err.Error())
	} else if err = DomainSeparator.Scan(hash); err != nil {
		return errors.New("config: domain separator load error: " + err.Error())
	}

	if ValidatorPkey, err = crypto.HexToECDSA(os.Getenv("VALIDATOR_PKEY")); err != nil {
		return errors.New("config: validator pkey load error: " + err.Error())
	} else if _, err = crypto.Sign(crypto.Keccak256(), ValidatorPkey); err != nil {
		return errors.New("config: try sign error: " + err.Error())
	}

	if ValidatorPort, err = strconv.ParseUint(os.Getenv("VALIDATOR_PORT"), 10, 16); err != nil {
		return errors.New("config: validator port load error: " + err.Error())
	}

	if MinStartDelay, err = strconv.ParseUint(os.Getenv("MIN_START_DELAY"), 10, 32); err != nil {
		return errors.New("config: min start delay load error: " + err.Error())
	}

	return nil
}

func InitExecutor() error {
	if err = Init(); err != nil {
		return err
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	if address, err := hex.DecodeString(os.Getenv("PGAS_ADDRESS")); err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	} else if err = PGasAddress.Scan(address); err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	}

	if address, err := hex.DecodeString(os.Getenv("TREASURY_ADDRESS")); err != nil {
		return errors.New("config: treasury address load error: " + err.Error())
	} else if err = TreasuryAddress.Scan(address); err != nil {
		return errors.New("config: treasury address load error: " + err.Error())
	}

	if ProviderWS, err = url.Parse(os.Getenv("PROVIDER_WS")); err != nil {
		return errors.New("config: provider load error: " + err.Error())
	}

	if ProviderHTTP, err = url.Parse(os.Getenv("PROVIDER_HTTP")); err != nil {
		return errors.New("config: provider load error: " + err.Error())
	}

	if ExecutorPkey, err = crypto.HexToECDSA(os.Getenv("EXECUTOR_PKEY")); err != nil {
		return errors.New("config: executor pkey load error: " + err.Error())
	} else if _, err = crypto.Sign(crypto.Keccak256(), ExecutorPkey); err != nil {
		return errors.New("config: try sign error: " + err.Error())
	}

	ExecutorAddress = crypto.PubkeyToAddress(ExecutorPkey.PublicKey)

	if ChainID, err = strconv.ParseUint(os.Getenv("CHAIN_ID"), 10, 64); err != nil {
		return errors.New("config: chain id load error: " + err.Error())
	}

	if PrevalidateDelay, err = strconv.ParseUint(os.Getenv("PREVALIDATE_DELAY"), 10, 32); err != nil {
		return errors.New("config: prevalidate delay load error: " + err.Error())
	}

	return nil
}
