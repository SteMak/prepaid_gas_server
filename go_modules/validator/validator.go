package validator

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

var (
	delay     uint32
	separator structs.Hash
	key       *ecdsa.PrivateKey

	err error
)

func Init(start_delay uint32, domain_separator structs.Hash, validator_key *ecdsa.PrivateKey) {
	delay = start_delay
	separator = domain_separator
	key = validator_key

	http.HandleFunc("/load", load)
	http.HandleFunc("/validate", validate)
}

func Start(port uint16) {
	for {
		if err = http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil); err != nil {
			log.Println("http: listen and serve error: " + err.Error())
		}
	}
}

func load(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	reverse, err := strconv.ParseBool(r.URL.Query().Get("reverse"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	messages, err := db.GetMessages(reverse, offset, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	data, err := json.Marshal(structs.WrapHTTPLoadResponses(messages))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, string(data))
}

func validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var request structs.HTTPValidateRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = utils.ValidateOffchain(request.Message, delay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	digest, err := request.Message.DigestHash(separator)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = digest.Verify(request.OrigSign, request.Message.From)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = utils.ValidateOnchain(request.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	valid_sign, err := digest.Sign(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	err = db.InsertMessage(structs.WrapDBMessage(request.Message, request.OrigSign, valid_sign))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, hex.EncodeToString(valid_sign[:]))
}
