package db

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/SteMak/prepaid_gas_server/go_modules/config"
	"github.com/SteMak/prepaid_gas_server/go_modules/structs"
)

var (
	DB *sql.DB

	err error
)

func Init() error {
	connect := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", config.PostgresUser, config.PostgresPassword)
	DB, err = sql.Open("postgres", connect)

	return err
}

func InsertMessage(message structs.Message, sign []byte, valid []byte) error {
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
		hex.EncodeToString(sign),
		hex.EncodeToString(valid),
	)

	return err
}
