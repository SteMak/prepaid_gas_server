package executor

import (
	"context"
	"encoding/hex"
	"math/big"

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

func acceptor(pgas_address common.Address) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{pgas_address},
		Topics: [][]common.Hash{{common.BytesToHash(crypto.Keccak256([]byte(
			"OrderCreate(uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)))",
		)))}},
	}

	events := make(chan types.Log)
	subscription, err := onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events)
	if err != nil {
		return err
	}

	for {
		select {
		case err = <-subscription.Err():
			for err != nil {
				subscription, err = onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events)
			}
		case event := <-events:
			order, err := onchain.WrapPGasOrder(event.Data)
			if err != nil {
				continue
			}

			id, err := structs.WrapUint256(event.Topics[1][:])
			if err != nil {
				continue
			}

			planOrder(id, order)
		}
	}
}

func planOrder(id structs.Uint256, order pgas.Order) {
	if orders[hex.EncodeToString(id[:])] != nil {
		return
	}

	if isOrderRisky(id, order) {
		return
	}

	orders[hex.EncodeToString(id[:])] = &order

	_, err = onchain.Treasury.OrderAccept(onchain.Transactor, id.ToBig())
	if err != nil {
		orders[hex.EncodeToString(id[:])] = nil
		return
	}
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
