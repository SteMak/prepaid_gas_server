package db

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

var (
	DB *sqlx.DB

	err error
)

func Init(user string, password string) error {
	connect := fmt.Sprintf(
		"user=%s password=%s dbname=postgres sslmode=disable",
		user, password,
	)
	if DB, err = sqlx.Connect("postgres", connect); err != nil {
		return errors.New("db: connection error: " + err.Error())
	}

	return nil
}

func InitValidator(user string, password string) error {
	if err = Init(user, password); err != nil {
		return err
	}

	if err = InitMessages(); err != nil {
		return err
	}

	return nil
}

func InitExecutor(user string, password string) error {
	if err = Init(user, password); err != nil {
		return err
	}

	if err = InitMessages(); err != nil {
		return err
	}

	return nil
}

func InitMessages() error {
	if sql, err := os.ReadFile(filepath.Join("sql", "messages.up.sql")); err != nil {
		return errors.New("db: messages up script read error: " + err.Error())
	} else if _, err = DB.Exec(string(sql)); err != nil {
		return errors.New("db: messages up script execution error: " + err.Error())
	}

	return nil
}

func GetMessages(reverse bool, offset uint64, limit uint64) ([]structs.DBMessage, error) {
	messages := []structs.DBMessage{}

	order := "asc"
	if reverse {
		order = "desc"
	}

	script := fmt.Sprintf(
		"select * from messages order by id %s limit %s offset %s",
		order,
		strconv.FormatUint(limit, 10),
		strconv.FormatUint(offset, 10),
	)
	err := DB.Select(&messages, script)

	return messages, err
}

func GetMessagesByOrder(order uint64, offset uint64, limit uint64) ([]structs.DBMessage, error) {
	messages := []structs.DBMessage{}

	script := fmt.Sprintf(
		"select * from messages where order_ = %s order by id limit %s offset %s",
		strconv.FormatUint(order, 10),
		strconv.FormatUint(limit, 10),
		strconv.FormatUint(offset, 10),
	)
	err := DB.Select(&messages, script)

	return messages, err
}

func InsertMessage(message structs.DBMessage) error {
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
			decode($9, 'hex')
		);`,
		hex.EncodeToString(message.From[:]),
		hex.EncodeToString(bytes.TrimLeft(message.Nonce[:], "\x00")),
		hex.EncodeToString(bytes.TrimLeft(message.Order[:], "\x00")),
		hex.EncodeToString(bytes.TrimLeft(message.Start[:], "\x00")),
		hex.EncodeToString(message.To[:]),
		hex.EncodeToString(bytes.TrimLeft(message.Gas[:], "\x00")),
		hex.EncodeToString(message.Data[:]),
		hex.EncodeToString(message.OrigSign[:]),
		hex.EncodeToString(message.ValidSign[:]),
	)

	return err
}
