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
	http.HandleFunc("/load", Load)
	http.HandleFunc("/validate", Validate)

	return http.ListenAndServe(":"+strconv.Itoa(config.ValidatorPort), nil)
}

func Load(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	reverse, err := strconv.ParseBool(r.URL.Query().Get("reverse"))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	messages, err := db.GetMessages(reverse, offset)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	data, err := json.Marshal(structs.WrapResponses(messages))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, string(data))
}

func Validate(w http.ResponseWriter, r *http.Request) {
	var request structs.HTTPValidateRequest

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

	err = request.OrigSign.Verify(digest, request.Message.Signer)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	valid, err := request.Message.Sign()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = db.InsertMessage(request.Message, request.OrigSign, valid)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, hex.EncodeToString(valid[:]))
}
