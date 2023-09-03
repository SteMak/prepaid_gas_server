package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/SteMak/prepaid_gas_server/config"
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

func InsertMessage() error {
	// TODO: Finalize
	_, err = DB.Exec("INSERT INTO messages(id) VALUES($1);", 7)
	return err
}
