package http

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/SteMak/prepaid_gas_server/go_modules/config"
	"github.com/SteMak/prepaid_gas_server/go_modules/db"
	"github.com/SteMak/prepaid_gas_server/go_modules/structs"
)

var (
	err error
)

func Init() error {
	// TODO: Add view endpoints
	http.HandleFunc("/", Validator)

	return http.ListenAndServe(":"+strconv.Itoa(config.ValidatorPort), nil)
}

func Validator(w http.ResponseWriter, r *http.Request) {
	var request structs.RequestInsert

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = request.Message.ValidateEarlyLiquidation(20)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	digest, err := request.Message.DigestHash()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = request.Sign.Verify(digest, request.Message.Signer)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	valid, err := request.Message.Sign()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = db.Insert(request.Message, request.Sign, valid)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, hex.EncodeToString(valid[:]))
}
