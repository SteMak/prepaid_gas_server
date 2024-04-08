package http

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/prepaidGas/prepaid-gas-server/go_modules/config"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/db"
	"github.com/prepaidGas/prepaid-gas-server/go_modules/structs"
)

var (
	err error
)

func Init() error {
	http.HandleFunc("/load", Load)
	http.HandleFunc("/validate", Validate)

	return http.ListenAndServe(":"+strconv.FormatUint(config.ValidatorPort, 10), nil)
}

func Load(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
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

	data, err := json.Marshal(structs.WrapHTTPLoadResponses(messages))
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

	err = request.Message.ValidateOffchain()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	digest, err := request.Message.DigestHash()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = request.OrigSign.Verify(digest, request.Message.From)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	valid_sign, err := digest.Sign()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	err = db.InsertMessage(structs.WrapDBMessage(request.Message, request.OrigSign, valid_sign))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, hex.EncodeToString(valid_sign[:]))
}
