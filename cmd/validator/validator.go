package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SteMak/prepaid_gas_server/structs"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

func SignMessage(message structs.Message) []byte {
	hash := structs.MessageHash(message)

	pkey, err := crypto.HexToECDSA(os.Getenv("VALIDATOR_PKEY"))
	if err != nil {
		fmt.Println("Bad PKEY:", err)
		return []byte{}
	}

	signature, err := crypto.Sign(hash, pkey)
	if err != nil {
		fmt.Println("Bad SIGN:", err)
		return []byte{}
	}

	return signature
}

func Validator(w http.ResponseWriter, r *http.Request) {
	var message structs.Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println("Bad MESS:", err)
		return
	}

	// TODO: Validate message fields has proper length

	if int64(message.Deadline) <= time.Now().Unix()+20 {
		io.WriteString(w, "Bad DEADLINE")
		return
	}

	io.WriteString(w, hex.EncodeToString(SignMessage(message)))
}

func main() {
	godotenv.Load()

	http.HandleFunc("/", Validator)

	err := http.ListenAndServe(":"+os.Getenv("VALIDATOR_PORT"), nil)
	if err != nil {
		fmt.Println("Bad HTTP:", err)
	}
}
