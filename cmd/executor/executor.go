package main

import (
	"log"

	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/executor"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
)

func main() {
	if err := config.InitExecutor(); err != nil {
		log.Fatalln(err.Error())
	}

	if err := onchain.InitExecutor(
		config.ProviderHTTP,
		config.ProviderWS,
		config.PGasAddress,
		config.TreasuryAddress,
		config.ExecutorPkey,
		config.GasFeeCap,
		config.GasTipCap,
		config.ChainID,
	); err != nil {
		log.Fatalln(err.Error())
	}

	if err := db.Init(config.PostgresUser, config.PostgresPassword, config.DBPort); err != nil {
		log.Fatalln(err.Error())
	}
	defer db.DB.Close()

	if err := executor.Init(config.PGasAddress, config.ExecutorAddress, config.PrevalidateDelay); err != nil {
		log.Fatalln(err.Error())
	}
	executor.Start()
}
