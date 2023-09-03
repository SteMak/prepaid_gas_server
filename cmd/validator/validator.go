package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"

	"github.com/SteMak/prepaid_gas_server/config"
	"github.com/SteMak/prepaid_gas_server/db"
	"github.com/SteMak/prepaid_gas_server/http"
)

var (
	err error
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = godotenv.Load()
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = config.Init()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = db.Init()
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = http.Init()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	wg.Wait()
}
