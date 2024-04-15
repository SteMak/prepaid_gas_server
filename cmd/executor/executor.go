package main

import (
	"log"

	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/executor"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
)

var (
	err error
)

func main() {
	err = config.InitExecutor()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = onchain.InitExecutor(
		config.ProviderHTTP, config.ProviderWS, config.PGasAddress, config.ExecutorPkey, config.ChainID,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = db.InitExecutor(config.PostgresUser, config.PostgresPassword)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.DB.Close()

	err = executor.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
