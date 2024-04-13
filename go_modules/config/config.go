package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

var (
	ProviderURL     *url.URL
	PGasAddress     structs.Address
	DomainSeparator structs.Hash

	ValidatorPkey *ecdsa.PrivateKey
	ValidatorPort uint64
	MinStartDelay uint64

	PostgresUser     string
	PostgresPassword string

	err error
)

func Init() error {
	ProviderURL, err = url.Parse(os.Getenv("PROVIDER_URL"))

	address, err := hex.DecodeString(os.Getenv("PGAS_ADDRESS"))
	PGasAddress.Scan(address)
	if err != nil {
		return errors.New("config: pgas address load error: " + err.Error())
	}

	hash, err := hex.DecodeString(os.Getenv("DOMAIN_SEPARATOR"))
	DomainSeparator.Scan(hash)
	if err != nil {
		return errors.New("config: domain separator load error: " + err.Error())
	}

	ValidatorPkey, err = crypto.HexToECDSA(os.Getenv("VALIDATOR_PKEY"))
	if err != nil {
		return errors.New("config: validator pkey load error: " + err.Error())
	}
	_, err = crypto.Sign(crypto.Keccak256(), ValidatorPkey)
	if err != nil {
		return errors.New("config: try sig failed: " + err.Error())
	}

	ValidatorPort, err = strconv.ParseUint(os.Getenv("VALIDATOR_PORT"), 10, 16)
	if err != nil {
		return errors.New("config: validator port load error: " + err.Error())
	}

	MinStartDelay, err = strconv.ParseUint(os.Getenv("MIN_START_DELAY"), 10, 32)
	if err != nil {
		return errors.New("config: min start delay load error: " + err.Error())
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	return nil
}
