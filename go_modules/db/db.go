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
)

func Init(user, password string, port uint16) error {
	connect := fmt.Sprintf(
		"user=%s password=%s dbname=postgres sslmode=disable port=%s",
		user, password, strconv.FormatUint(uint64(port), 10),
	)
	if db, err := sqlx.Connect("postgres", connect); err != nil {
		return errors.New("db: connection error: " + err.Error())
	} else {
		DB = db
	}

	if err := InitMessages(); err != nil {
		return err
	}

	return nil
}

func InitMessages() error {
	if sql, err := os.ReadFile(filepath.Join("sql", "messages.up.sql")); err != nil {
		return errors.New("db: messages up script read error: " + err.Error())
	} else if _, err := DB.Exec(string(sql)); err != nil {
		return errors.New("db: messages up script execution error: " + err.Error())
	}

	return nil
}

func GetMessages(reverse bool, offset, limit uint64) ([]structs.DBMessage, error) {
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

	if err := DB.Select(&messages, script); err != nil {
		return messages, errors.New("db: get messages: " + err.Error())
	}

	return messages, nil
}

func GetMessagesByOrder(order structs.Uint256, offset, limit uint64) ([]structs.DBMessage, error) {
	messages := []structs.DBMessage{}

	script := fmt.Sprintf(
		"select * from messages where order_ = decode('%s', 'hex') order by id limit %s offset %s",
		hex.EncodeToString(order[:]),
		strconv.FormatUint(limit, 10),
		strconv.FormatUint(offset, 10),
	)

	if err := DB.Select(&messages, script); err != nil {
		return messages, errors.New("db: get messages by order: " + err.Error())
	}

	return messages, nil
}

func InsertMessage(message structs.DBMessage) error {
	_, err := DB.Exec(`insert into messages
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
	if err != nil {
		return errors.New("db: insert message: " + err.Error())
	}

	return nil
}
