package utils

import (
	"errors"
	"time"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func ValidateOffchain(message structs.Message, min_delay uint32) error {
	if start, err := message.Start.ToUint32(); err != nil {
		return errors.New("message: message start parse: " + err.Error())
	} else if int64(start) <= time.Now().Unix()+int64(min_delay) {
		return errors.New("message: message provided lately")
	}

	if err = message.From.NotZero(); err != nil {
		return errors.New("message: message from zero: " + err.Error())
	}

	return nil
}

func ValidateOnchain(message structs.Message) error {
	// Error handling omitted to not stuck due to node errors
	result, _ := onchain.PGas.MessageValidate(nil, onchain.WrapPGasMessage(message))

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
