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
		if port == 443 {
			if err := http.ListenAndServeTLS(
				":"+strconv.FormatUint(uint64(port), 10), "ssl/server.crt", "ssl/server.key", nil,
			); err != nil {
				log.Printf("listen and serve https: %s\n\n", err.Error())
			}
		} else {
			if err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil); err != nil {
				log.Printf("listen and serve http: %s\n\n", err.Error())
			}
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

		log.Printf("\"%s\" parsing offset: %s\n\n", r.URL.String(), err.Error())
		return
	}

	reverse, err := strconv.ParseBool(r.URL.Query().Get("reverse"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "parsing reverse: "+err.Error())

		log.Printf("\"%s\" parsing reverse: %s\n\n", r.URL.String(), err.Error())
		return
	}

	messages, err := db.GetMessages(reverse, offset, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "db query: "+err.Error())

		log.Printf("\"%s\" db query: %s\n\n", r.URL.String(), err.Error())
		return
	}

	data, err := json.Marshal(structs.WrapHTTPLoadResponses(messages))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "reply marchal: "+err.Error())

		log.Printf("\"%s\" reply marchal: %s\n\n", r.URL.String(), err.Error())
		return
	}

	io.WriteString(w, string(data))
}

func validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var request structs.HTTPValidateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "parsing body: "+err.Error())

		log.Printf("\"%s\" parsing body: %s\n\n", r.URL.String(), err)
		return
	}

	if err := utils.ValidateOffchain(request.Message, delay); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "offchain validation: "+err.Error())

		log.Printf("\"%s\" offchain validation: \"%#v\": %s\n\n", r.URL.String(), request, err)
		return
	}

	digest, err := request.Message.DigestHash(separator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "digest hash: "+err.Error())

		log.Printf("\"%s\" digest hash: \"%#v\": %s\n\n", r.URL.String(), request, err)
		return
	}

	if err := digest.Verify(request.OrigSign, request.Message.From); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "bad signature: "+err.Error())

		log.Printf("\"%s\" bad signature: \"%#v\": %s\n\n", r.URL.String(), request, err)
		return
	}

	if err := utils.ValidateOnchain(request.Message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "onchain validation: "+err.Error())

		log.Printf("\"%s\" onchain validation: \"%#v\": %s\n\n", r.URL.String(), request, err)
		return
	}

	valid_sign, err := digest.Sign(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "reply signature: "+err.Error())

		log.Printf("\"%s\" reply signature: \"%#v\": %s\n\n", r.URL.String(), request, err)
		return
	}

	db_message := structs.WrapDBMessage(request.Message, request.OrigSign, valid_sign)
	if err := db.InsertMessage(db_message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "db insert: "+err.Error())

		log.Printf("\"%s\" db insert: \"%#v\": %s\n\n", r.URL.String(), db_message, err)
		return
	}

	io.WriteString(w, "0x"+hex.EncodeToString(valid_sign[:]))
	log.Printf("\"%s\" success: \"%#v\"\n\n", r.URL.String(), db_message)
}
