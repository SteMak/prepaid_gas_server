package main

import (
	"log"

	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/validator"
)

var (
	err error
)

func main() {
	err = config.InitValidator()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = onchain.InitValidator(config.ProviderHTTP, config.PGasAddress, config.DomainSeparator)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = db.InitValidator(config.PostgresUser, config.PostgresPassword)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.DB.Close()

	err = validator.Init(config.ValidatorPort)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
