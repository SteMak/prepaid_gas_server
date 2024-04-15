package executor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

var (
	orders       map[*big.Int]*pgas.Order
	last_message uint64

	err error
)

func Init() error {
	{
		number, err := onchain.ClientHTTP.BlockNumber(context.Background())
		if err != nil {
			return err
		}
		offset := int64(0)
		for {
			result, err := onchain.PGas.GetExecutorOrders(
				&bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(number)},
				config.ExecutorAddress, true, big.NewInt(100), big.NewInt(offset),
			)
			if err != nil {
				return err
			}
			for _, item := range result {
				orders[item.Id] = &item.Order
			}
			if len(result) < 100 {
				break
			}
			offset += 100
		}
	}
	{
		offset := uint64(0)
		for {
			result, err := db.GetMessages(false, offset, 100)
			if err != nil {
				return err
			}
			for _, item := range result {
				if orders[big.NewInt(0).SetBytes(item.Order[:])] != nil {
					message, sign, _ := structs.UnwrapDBMessage(item)
					go PlanMessage(message, sign)
				}
			}
			last_message += uint64(len(result))
			if len(result) < 100 {
				break
			}
			offset += 100
		}
	}

	return nil
}

func MonitorMessages() {}

func RunMessage(message structs.Message, sign structs.Signature) error {
	tx, err := onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign[:])
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())

	return nil
}

func PlanMessage(message structs.Message, sign structs.Signature) error {
	order := orders[big.NewInt(0).SetBytes(message.Order[:])]
	if order == nil {
		return nil
	}

	start, _ := message.Start.ToUint32()
	ex_start := int64(start)
	ex_end := int64(ex_start) + order.TxWindow.Int64()

	// if ex_start < order.Start.Int64() || order.End.Int64() < ex_end {
	// 	return nil
	// }

	if ex_end < time.Now().Unix() {
		return nil
	}

	if ex_start-time.Now().Unix() > 100 {
		time.Sleep(time.Duration(ex_start - 100 - time.Now().Unix()))
		used, _ := onchain.PGas.Nonce(nil, common.Address(message.From), big.NewInt(0).SetBytes(message.Nonce[:]))
		if used {
			return nil
		}
	}

	time.Sleep(time.Duration(ex_start - time.Now().Unix()))

	return RunMessage(message, sign)
}
