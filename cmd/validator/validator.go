package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/prepaidGas/prepaid-gas-server/go_modules/config"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/db"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/http"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/onchain"
)

var (
	err error
)

func main() {
	err = godotenv.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = config.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = onchain.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = db.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.DB.Close()

	err = db.InitMessages()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = http.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
