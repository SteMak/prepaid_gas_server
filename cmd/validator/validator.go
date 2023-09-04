package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/SteMak/prepaid_gas_server/go_modules/config"
	"github.com/SteMak/prepaid_gas_server/go_modules/db"
	"github.com/SteMak/prepaid_gas_server/go_modules/http"
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

	err = db.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = http.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
