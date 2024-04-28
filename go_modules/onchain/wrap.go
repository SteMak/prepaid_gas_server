package onchain

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func WrapPGasMessage(message structs.Message) pgas.Message {
	return pgas.Message{
		From:  common.Address(message.From),
		Nonce: message.Nonce.ToBig(),
		Order: message.Order.ToBig(),
		Start: message.Start.ToBig(),
		To:    common.Address(message.To),
		Gas:   message.Gas.ToBig(),
		Data:  message.Data,
	}
}

func WrapPGasOrder(data []byte) (pgas.Order, error) {
	if len(data) != 352 {
		return pgas.Order{}, errors.New("onchain: incorrect data length")
	}

	return pgas.Order{
		Manager:      common.BytesToAddress(data[0:32]),
		Gas:          big.NewInt(0).SetBytes(data[32:64]),
		Expire:       big.NewInt(0).SetBytes(data[64:96]),
		Start:        big.NewInt(0).SetBytes(data[96:128]),
		End:          big.NewInt(0).SetBytes(data[128:160]),
		TxWindow:     big.NewInt(0).SetBytes(data[160:192]),
		RedeemWindow: big.NewInt(0).SetBytes(data[192:224]),
		GasPrice: pgas.GasPayment{
			Token:   common.BytesToAddress(data[224:256]),
			PerUnit: big.NewInt(0).SetBytes(data[256:288]),
		},
		GasGuarantee: pgas.GasPayment{
			Token:   common.BytesToAddress(data[288:320]),
			PerUnit: big.NewInt(0).SetBytes(data[320:352]),
		},
	}, nil
}
