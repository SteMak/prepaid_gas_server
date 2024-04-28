package main

import (
	"log"

	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/validator"
)

func main() {
	if err := config.InitValidator(); err != nil {
		log.Fatalln(err.Error())
	}

	if err := onchain.InitValidator(config.ProviderHTTP, config.PGasAddress, config.DomainSeparator); err != nil {
		log.Fatalln(err.Error())
	}

	if err := db.Init(config.PostgresUser, config.PostgresPassword); err != nil {
		log.Fatalln(err.Error())
	}
	defer db.DB.Close()

	validator.Init(config.MinStartDelay, config.DomainSeparator, config.ValidatorPkey)
	validator.Start(config.ValidatorPort)
}
