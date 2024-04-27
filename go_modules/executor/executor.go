package executor

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/config"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

var (
	orders = make(map[string]*pgas.Order)
	offset uint64

	err error
)

func Init() error {
	if err = FillOrders(); err != nil {
		return err
	}
	if err = FillMessages(); err != nil {
		return err
	}

	go MonitorMessages()
	Acceptor(config.PGasAddress)

	return nil
}

func FillOrders() error {
	number, err := onchain.ClientHTTP.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	opts := &bind.CallOpts{BlockNumber: big.NewInt(0).SetUint64(number)}

	limit := int64(100)
	offset := int64(0)
	for {
		if result, err := onchain.PGas.GetExecutorOrders(
			opts, config.ExecutorAddress, true, big.NewInt(limit), big.NewInt(offset),
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
			go PlanMessage(message, sign)
		}

		offset += uint64(len(result))
		if uint64(len(result)) < limit {
			break
		}
	}

	return nil
}

func MonitorMessages() {
	for {
		result, _ := db.GetMessages(false, offset, 1000)
		offset += uint64(len(result))

		for _, item := range result {
			message, sign, _ := structs.UnwrapDBMessage(item)
			go PlanMessage(message, sign)
		}

		time.Sleep(time.Second * 5)
	}
}

func RunMessage(message structs.Message, sign structs.Signature) error {
	_, err := onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
	if err != nil {
		// TODO: Make optional for local node
		transactor := *onchain.Transactor
		transactor.GasFeeCap = big.NewInt(1000)
		transactor.GasTipCap = big.NewInt(1000)

		_, err = onchain.PGas.Execute(&transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
		return err
	}

	return nil
}

func PlanOrder(id structs.Uint256, order pgas.Order) {
	if orders[hex.EncodeToString(id[:])] != nil {
		return
	}

	orders[hex.EncodeToString(id[:])] = &order

	messages, err := db.GetMessagesByOrder(id, 0, 1)
	if err != nil || uint64(len(messages)) > 0 {
		return
	}

	if order.GasGuarantee.PerUnit.Cmp(big.NewInt(0)) != 0 {
		return
	}

	_, err = onchain.Treasury.OrderAccept(onchain.Transactor, id.ToBig())
	if err != nil {
		// TODO: Make optional for local node
		transactor := *onchain.Transactor
		transactor.GasFeeCap = big.NewInt(1000)
		transactor.GasTipCap = big.NewInt(1000)

		_, err = onchain.Treasury.OrderAccept(&transactor, id.ToBig())
		if err != nil {
			orders[hex.EncodeToString(id[:])] = nil
			return
		}
	}
}

func PlanMessage(message structs.Message, sign structs.Signature) {
	order := orders[hex.EncodeToString(message.Order[:])]
	if order == nil {
		return
	}

	// start + window < now
	if big.NewInt(0).Add(message.Start.ToBig(), order.TxWindow).Cmp(utils.UnixBig()) == -1 {
		return
	}

	// start - now > delay
	if big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()).Cmp(big.NewInt(int64(config.PrevalidateDelay))) == 1 {
		time.Sleep(time.Second * time.Duration(
			big.NewInt(0).Sub(
				big.NewInt(0).Sub(
					message.Start.ToBig(), utils.UnixBig(),
				),
				big.NewInt(int64(config.PrevalidateDelay)),
			).Int64(),
		))
		used, _ := onchain.PGas.Nonce(nil, common.Address(message.From), message.Nonce.ToBig())
		if used {
			return
		}
	}

	time.Sleep(time.Second * time.Duration(big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()).Int64()))

	_ = RunMessage(message, sign)
}
