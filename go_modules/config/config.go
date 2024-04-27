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

	PGasAddress     common.Address
	TreasuryAddress common.Address

	GasFeeCap       *int64
	GasTipCap       *int64
	ChainID         uint64
	DomainSeparator structs.Hash

	ValidatorPkey   *ecdsa.PrivateKey
	ExecutorPkey    *ecdsa.PrivateKey
	ExecutorAddress common.Address

	MinStartDelay    uint32
	PrevalidateDelay uint32

	ValidatorPort uint16

	err error
)

func InitValidator() error {
	if err = loadEnv(); err != nil {
		return err
	}

	loadPostgres()
	if err = loadProvider(true, false); err != nil {
		return err
	}
	if err = loadAddresses(true, false); err != nil {
		return err
	}
	if err = loadChainDetails(false, false, false, true); err != nil {
		return err
	}
	if err = loadPkey(true, false); err != nil {
		return err
	}
	if err = loadDelays(true, false); err != nil {
		return err
	}
	if err = loadHTTP(true); err != nil {
		return err
	}

	return nil
}

func InitExecutor() error {
	if err = loadEnv(); err != nil {
		return err
	}

	loadPostgres()
	if err = loadProvider(true, true); err != nil {
		return err
	}
	if err = loadAddresses(true, true); err != nil {
		return err
	}
	if err = loadChainDetails(true, true, true, false); err != nil {
		return err
	}
	if err = loadPkey(false, true); err != nil {
		return err
	}
	if err = loadDelays(false, true); err != nil {
		return err
	}
	if err = loadHTTP(false); err != nil {
		return err
	}

	return nil
}

func loadEnv() error {
	if err = godotenv.Load(); err != nil {
		return errors.New("config: environment load error: " + err.Error())
	}

	return nil
}

func loadPostgres() {
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
}

func loadProvider(http, websocket bool) error {
	if ProviderHTTP, err = url.Parse(os.Getenv("PROVIDER_HTTP")); http && err != nil {
		return errors.New("config: provider load error: " + err.Error())
	}

	if ProviderWS, err = url.Parse(os.Getenv("PROVIDER_WS")); websocket && err != nil {
		return errors.New("config: provider load error: " + err.Error())
	}

	return nil
}

func loadAddresses(pgas, treasury bool) error {
	if address, err := hex.DecodeString(os.Getenv("PGAS_ADDRESS")); pgas && err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	} else {
		PGasAddress = common.BytesToAddress(address)
	}

	if address, err := hex.DecodeString(os.Getenv("TREASURY_ADDRESS")); treasury && err != nil {
		return errors.New("config: treasury address load error: " + err.Error())
	} else {
		TreasuryAddress = common.BytesToAddress(address)
	}

	return nil
}

func loadChainDetails(gasfeecap, gastipcap, chain_id, separator bool) error {
	if num, err := strconv.ParseInt(os.Getenv("GAS_FEE_CAP"), 10, 64); gasfeecap && err != nil {
		return errors.New("config: gas fee cap load error: " + err.Error())
	} else if gasfeecap && num >= 0 {
		GasFeeCap = &num
	}

	if num, err := strconv.ParseInt(os.Getenv("GAS_TIP_CAP"), 10, 64); gastipcap && err != nil {
		return errors.New("config: gas tip cap load error: " + err.Error())
	} else if gastipcap && num >= 0 {
		GasTipCap = &num
	}

	if ChainID, err = strconv.ParseUint(os.Getenv("CHAIN_ID"), 10, 64); chain_id && err != nil {
		return errors.New("config: chain id load error: " + err.Error())
	}

	if hash, err := hex.DecodeString(os.Getenv("DOMAIN_SEPARATOR")); separator && err != nil {
		return errors.New("config: domain separator load error: " + err.Error())
	} else if err = DomainSeparator.Scan(hash); separator && err != nil {
		return errors.New("config: domain separator load error: " + err.Error())
	}

	return nil
}

func loadPkey(validator, executor bool) error {
	if ValidatorPkey, err = crypto.HexToECDSA(os.Getenv("VALIDATOR_PKEY")); validator && err != nil {
		return errors.New("config: validator pkey load error: " + err.Error())
	} else if _, err = crypto.Sign(crypto.Keccak256(), ValidatorPkey); validator && err != nil {
		return errors.New("config: try sign error: " + err.Error())
	}

	if ExecutorPkey, err = crypto.HexToECDSA(os.Getenv("EXECUTOR_PKEY")); executor && err != nil {
		return errors.New("config: executor pkey load error: " + err.Error())
	} else if _, err = crypto.Sign(crypto.Keccak256(), ExecutorPkey); executor && err != nil {
		return errors.New("config: try sign error: " + err.Error())
	}

	ExecutorAddress = crypto.PubkeyToAddress(ExecutorPkey.PublicKey)

	return nil
}

func loadDelays(v_start_delay, x_check_delay bool) error {
	if num, err := strconv.ParseUint(os.Getenv("MIN_START_DELAY"), 10, 32); v_start_delay && err != nil {
		return errors.New("config: min start delay load error: " + err.Error())
	} else {
		MinStartDelay = uint32(num)
	}

	if num, err := strconv.ParseUint(os.Getenv("PREVALIDATE_DELAY"), 10, 32); x_check_delay && err != nil {
		return errors.New("config: prevalidate delay load error: " + err.Error())
	} else {
		PrevalidateDelay = uint32(num)
	}

	return nil
}

func loadHTTP(validator_port bool) error {
	if num, err := strconv.ParseUint(os.Getenv("VALIDATOR_PORT"), 10, 16); validator_port && err != nil {
		return errors.New("config: validator port load error: " + err.Error())
	} else {
		ValidatorPort = uint16(num)
	}

	return nil
}
