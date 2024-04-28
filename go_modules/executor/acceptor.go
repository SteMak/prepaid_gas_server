package executor

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

var (
	query        ethereum.FilterQuery
	events       = make(chan types.Log)
	subscription ethereum.Subscription
)

func initAcceptor(pgas_address common.Address) error {
	query = ethereum.FilterQuery{
		Addresses: []common.Address{pgas_address},
		Topics: [][]common.Hash{{common.BytesToHash(crypto.Keccak256([]byte(
			"OrderCreate(uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)))",
		)))}},
	}

	if sub, err := onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events); err != nil {
		return errors.New("subscription: " + err.Error())
	} else {
		subscription = sub
	}

	return nil
}

func acceptor() {
	for {
		select {
		case err := <-subscription.Err():
			for err != nil {
				log.Printf("subscription: %s\n", err.Error())
				time.Sleep(time.Second)

				subscription, err = onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events)
			}
		case event := <-events:
			order, err := onchain.WrapPGasOrder(event.Data)
			if err != nil {
				log.Printf("event data not order: \"%#v\": %s\n", event, err.Error())
				continue
			}

			id, err := structs.WrapUint256(event.Topics[1][:])
			if err != nil {
				log.Printf("event topic not id: \"%#v\": %s\n", event, err.Error())
				continue
			}

			planOrder(id, order)
		}
	}
}

func planOrder(id structs.Uint256, order pgas.Order) {
	if orders[id.ToString()] != nil {
		log.Printf("order exists: %s\n", id.ToString())
		return
	}

	orders[id.ToString()] = &order

	if isOrderRisky(id, order) {
		log.Printf("order risky: %s\n", id.ToString())
		return
	}

	_, err := onchain.Treasury.OrderAccept(onchain.Transactor, id.ToBig())
	if err != nil {
		orders[id.ToString()] = nil

		log.Printf("order accept: %s: %s\n", id.ToString(), err.Error())
		return
	}

	log.Printf("order accept success: %s\n", id.ToString())
}

func isOrderRisky(id structs.Uint256, order pgas.Order) bool {
	messages, err := db.GetMessagesByOrder(id, 0, 1)
	if err != nil || uint64(len(messages)) > 0 {
		return true
	}

	if order.GasGuarantee.PerUnit.Cmp(big.NewInt(0)) == 0 {
		return false
	}

	if order.GasGuarantee.PerUnit.Cmp(big.NewInt(100001)) == -1 &&
		order.Gas.Cmp(big.NewInt(1000001)) == -1 &&
		big.NewInt(0).Add(utils.UnixBig(), big.NewInt(60*60*30)).Cmp(order.End) == 1 {
		return false
	}

	return true
}
