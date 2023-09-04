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
	ValidatorPort   int
	DomainSeparator []byte

	PostgresUser     string
	PostgresPassword string

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

	ValidatorPort, err = strconv.Atoi(os.Getenv("VALIDATOR_PORT"))
	if err != nil {
		return errors.New("config: validator port load error")
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	return nil
}
