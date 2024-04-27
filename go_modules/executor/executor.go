package executor

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

var (
	orders = make(map[string]*pgas.Order)
	offset uint64
	delay  uint32

	err error
)

func Init(executor common.Address, prevalidate_delay uint32) error {
	if err = FillOrders(executor); err != nil {
		return err
	}
	if err = FillMessages(); err != nil {
		return err
	}

	delay = prevalidate_delay

	return nil
}

func Start(pgas_address common.Address) {
	go monitorMessages()
	acceptor(pgas_address)
}

func FillOrders(executor common.Address) error {
	number, err := onchain.ClientHTTP.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	opts := &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(number)}

	limit := int64(100)
	offset := int64(0)
	for {
		if result, err := onchain.PGas.GetExecutorOrders(
			opts, executor, true, big.NewInt(limit), big.NewInt(offset),
		); err != nil {
			return err
		} else {
			for _, item := range result {
				var id structs.Uint256
				id.Scan(item.Id.Bytes())
				orders[hex.EncodeToString(id[:])] = &item.Order
			}

			if int64(len(result)) < limit {
				break
			}
		}

		offset += limit
	}

	return nil
}

func FillMessages() error {
	offset = uint64(0)
	limit := uint64(100)
	for {
		result, err := db.GetMessages(false, limit, 100)
		if err != nil {
			return err
		}

		for _, item := range result {
			message, sign, _ := structs.UnwrapDBMessage(item)
			go planMessage(message, sign)
		}

		offset += uint64(len(result))
		if uint64(len(result)) < limit {
			break
		}
	}

	return nil
}
