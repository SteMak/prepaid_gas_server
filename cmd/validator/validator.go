package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/SteMak/prepaid_gas_server/config"
	"github.com/SteMak/prepaid_gas_server/structs"
	"github.com/joho/godotenv"
)

var (
	err error
)

func Validator(w http.ResponseWriter, r *http.Request) {
	var message structs.Message

	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = message.ValidateEarlyLiquidation(20)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	sign, err := structs.SignMessage(message)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, hex.EncodeToString(sign))
}

func main() {
	err = godotenv.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = config.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	http.HandleFunc("/", Validator)

	err = http.ListenAndServe(":"+strconv.Itoa(config.ValidatorPort), nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
