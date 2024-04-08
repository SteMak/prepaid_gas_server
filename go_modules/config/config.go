package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ValidatorPkey   *ecdsa.PrivateKey
	ValidatorPort   uint64
	DomainSeparator []byte

	PostgresUser     string
	PostgresPassword string

	MinStartDelay uint64

	err error
)

func Init() error {
	DomainSeparator, err = hex.DecodeString(os.Getenv("DOMAIN_SEPARATOR"))
	if err != nil {
		return errors.New("config: domain separator load error")
	}

	ValidatorPkey, err = crypto.HexToECDSA(os.Getenv("VALIDATOR_PKEY"))
	if err != nil {
		return errors.New("config: validator pkey load error")
	}
	_, err := crypto.Sign(crypto.Keccak256(), ValidatorPkey)
	if err != nil {
		return errors.New("config: try sig failed: " + err.Error())
	}

	ValidatorPort, err = strconv.ParseUint(os.Getenv("VALIDATOR_PORT"), 10, 16)
	if err != nil {
		return errors.New("config: validator port load error")
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	MinStartDelay, err = strconv.ParseUint(os.Getenv("MIN_START_DELAY"), 10, 32)
	if err != nil {
		return errors.New("config: min start delay load error")
	}

	return nil
}
