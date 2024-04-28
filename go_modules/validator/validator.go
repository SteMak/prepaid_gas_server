package validator

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
			log.Printf("listen and serve: %s\n", err.Error())
		}

		time.Sleep(time.Second)
	}
}

func load(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "parsing offset: "+err.Error())

		log.Printf("\"%s\" parsing offset: %s\n", r.URL.String(), err.Error())
		return
	}

	reverse, err := strconv.ParseBool(r.URL.Query().Get("reverse"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "parsing reverse: "+err.Error())

		log.Printf("\"%s\" parsing reverse: %s\n", r.URL.String(), err.Error())
		return
	}

	messages, err := db.GetMessages(reverse, offset, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "db query: "+err.Error())

		log.Printf("\"%s\" db query: %s\n", r.URL.String(), err.Error())
		return
	}

	data, err := json.Marshal(structs.WrapHTTPLoadResponses(messages))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "reply marchal: "+err.Error())

		log.Printf("\"%s\" reply marchal: %s\n", r.URL.String(), err.Error())
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
		io.WriteString(w, "parsing body: "+err.Error())

		log.Printf("\"%s\" parsing body: %s\n", r.URL.String(), err)
		return
	}

	err = utils.ValidateOffchain(request.Message, delay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "offchain validation: "+err.Error())

		log.Printf("\"%s\" offchain validation: \"%#v\": %s\n", r.URL.String(), request, err)
		return
	}

	digest, err := request.Message.DigestHash(separator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "digest hash: "+err.Error())

		log.Printf("\"%s\" digest hash: \"%#v\": %s\n", r.URL.String(), request, err)
		return
	}

	err = digest.Verify(request.OrigSign, request.Message.From)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "bad signature: "+err.Error())

		log.Printf("\"%s\" bad signature: \"%#v\": %s\n", r.URL.String(), request, err)
		return
	}

	err = utils.ValidateOnchain(request.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "onchain validation: "+err.Error())

		log.Printf("\"%s\" onchain validation: \"%#v\": %s\n", r.URL.String(), request, err)
		return
	}

	valid_sign, err := digest.Sign(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "reply signature: "+err.Error())

		log.Printf("\"%s\" reply signature: \"%#v\": %s\n", r.URL.String(), request, err)
		return
	}

	db_message := structs.WrapDBMessage(request.Message, request.OrigSign, valid_sign)
	err = db.InsertMessage(db_message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "db insert: "+err.Error())

		log.Printf("\"%s\" db insert: \"%#v\": %s\n", r.URL.String(), db_message, err)
		return
	}

	io.WriteString(w, hex.EncodeToString(valid_sign[:]))
	log.Printf("\"%s\" success: \"%#v\"\n", r.URL.String(), db_message)
}
