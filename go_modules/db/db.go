package db

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/SteMak/prepaid_gas_server/go_modules/config"
	"github.com/SteMak/prepaid_gas_server/go_modules/structs"
)

var (
	DB *sqlx.DB

	err error
)

func Init() error {
	connect := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", config.PostgresUser, config.PostgresPassword)
	DB, err = sqlx.Connect("postgres", connect)

	return err
}

func GetMessages(reverse bool, offset int) ([]structs.DBMessage, error) {
	messages := []structs.DBMessage{}

	order := "asc"
	if reverse {
		order = " desc"
	}

	script := fmt.Sprintf("select * from messages order by id %s limit 1000 offset %s", order, strconv.Itoa(offset))
	err := DB.Select(&messages, script)

	return messages, err
}

func InsertMessage(message structs.Message, orig_sign structs.Signature, valid_sign structs.Signature) error {
	_, err = DB.Exec(`insert into messages
		values(
			decode($1, 'hex'),
			decode($2, 'hex'),
			decode($3, 'hex'),
			decode($4, 'hex'),
			decode($5, 'hex'),
			decode($6, 'hex'),
			decode($7, 'hex'),
			decode($8, 'hex'),
			decode($9, 'hex'),
			decode($10, 'hex')
		);`,
		hex.EncodeToString(message.Signer[:]),
		hex.EncodeToString(bytes.TrimLeft(message.Nonce[:], "\x00")),
		hex.EncodeToString(bytes.TrimLeft(message.GasOrder[:], "\x00")),
		hex.EncodeToString(message.OnBehalf[:]),
		hex.EncodeToString(bytes.TrimLeft(message.Deadline[:], "\x00")),
		hex.EncodeToString(bytes.TrimLeft(message.Endpoint[:], "\x00")),
		hex.EncodeToString(bytes.TrimLeft(message.Gas[:], "\x00")),
		hex.EncodeToString(message.Data[:]),
		hex.EncodeToString(orig_sign[:]),
		hex.EncodeToString(valid_sign[:]),
	)

	return err
}
