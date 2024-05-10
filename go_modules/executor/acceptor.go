package executor

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

var (
	query        ethereum.FilterQuery
	events       = make(chan types.Log)
	subscription ethereum.Subscription
	sub_renew    *time.Ticker
)

func initAcceptor(pgas_address common.Address, renew_time uint32) error {
	query = ethereum.FilterQuery{
		Addresses: []common.Address{pgas_address},
		Topics: [][]common.Hash{{common.BytesToHash(crypto.Keccak256([]byte(
			"OrderCreate(uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)))",
		)))}},
	}

	sub_renew = time.NewTicker(time.Duration(renew_time) * time.Second)

	if err := subscribe(); err != nil {
		return err
	}

	return nil
}

func acceptor() {
	for {
		select {
		case <-sub_renew.C:
			log.Printf("subscription: timer renew\n\n")
			resubscribe()
		case err := <-subscription.Err():
			if err != nil {
				log.Printf("subscription: %s\n\n", err.Error())
			} else {
				log.Printf("subscription: dead nil\n\n")
			}
			resubscribe()
		case event := <-events:
			order, err := onchain.WrapPGasOrder(event.Data)
			if err != nil {
				log.Printf("event data not order: \"%#v\": %s\n\n", event, err.Error())
				continue
			}

			id, err := structs.WrapUint256(event.Topics[1][:])
			if err != nil {
				log.Printf("event topic not id: \"%#v\": %s\n\n", event, err.Error())
				continue
			}

			planOrder(id, order)
		}
	}
}

func subscribe() error {
	if sub, err := onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events); err != nil {
		return errors.New("subscribe: " + err.Error())
	} else {
		subscription = sub
	}

	return nil
}

func resubscribe() {
	err := subscribe()
	for err != nil {
		log.Printf("re%s\n\n", err.Error())
		time.Sleep(time.Second)

		err = subscribe()
	}
}

func planOrder(id structs.Uint256, order pgas.Order) {
	if orders[id.ToString()] != nil {
		log.Printf("order exists: %s\n\n", id.ToString())
		return
	}

	orders[id.ToString()] = &order

	if utils.IsOrderRisky(id, order) {
		log.Printf("order risky: %s\n\n", id.ToString())
		return
	}

	_, err := onchain.Treasury.OrderAccept(onchain.Transactor, id.ToBig())
	if err != nil {
		orders[id.ToString()] = nil

		log.Printf("order accept: %s: %s\n\n", id.ToString(), err.Error())
		return
	}

	log.Printf("order accept success: %s\n\n", id.ToString())
}
