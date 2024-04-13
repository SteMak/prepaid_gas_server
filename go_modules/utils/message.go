package utils

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func ValidateOffchain(message structs.Message, min_delay uint64) error {
	start, err := message.Start.ToUint32()
	if err != nil {
		return err
	}

	if int64(start) <= time.Now().Unix()+int64(min_delay) {
		return errors.New("message: message provided lately")
	}

	err = message.From.NotZero()

	return err
}

func ValidateOnchain(message structs.Message) error {
	// Error handling omitted to not stuck due to node errors
	result, _ := onchain.PGas.MessageValidate(nil, pgas.Message{
		From:  common.Address(message.From),
		Nonce: big.NewInt(0).SetBytes(message.Nonce[:]),
		Order: big.NewInt(0).SetBytes(message.Order[:]),
		Start: big.NewInt(0).SetBytes(message.Start[:]),
		To:    common.Address(message.To),
		Gas:   big.NewInt(0).SetBytes(message.Gas[:]),
		Data:  message.Data,
	})

	switch onchain.Validation(result) {
	case onchain.StartInFuture:
		return errors.New("onchain: start not in future")
	case onchain.NonceExhaustion:
		return errors.New("onchain: nonce already used")
	case onchain.BalanceCompliance:
		return errors.New("onchain: low order balance")
	case onchain.OwnerCompliance:
		return errors.New("onchain: incompliant order owner")
	case onchain.TimelineCompliance:
		return errors.New("onchain: not in order timeline")
	}

	return nil
}
