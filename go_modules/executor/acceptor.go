package executor

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func Acceptor(pgas_address structs.Address) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.BytesToAddress(pgas_address[:])},
		Topics:    [][]common.Hash{{common.BytesToHash(crypto.Keccak256([]byte("OrderCreate(uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)))")))}},
	}

	events := make(chan types.Log)
	subscription, err := onchain.ClientWS.SubscribeFilterLogs(context.Background(), query, events)
	if err != nil {
		return err
	}

	for {
		select {
		case err = <-subscription.Err():
			return err
		case event := <-events:
			order := pgas.Order{
				Manager:      common.BytesToAddress(event.Data[0:32]),
				Gas:          big.NewInt(0).SetBytes(event.Data[32:64]),
				Expire:       big.NewInt(0).SetBytes(event.Data[64:96]),
				Start:        big.NewInt(0).SetBytes(event.Data[96:128]),
				End:          big.NewInt(0).SetBytes(event.Data[128:160]),
				TxWindow:     big.NewInt(0).SetBytes(event.Data[160:192]),
				RedeemWindow: big.NewInt(0).SetBytes(event.Data[192:224]),
				GasPrice: pgas.GasPayment{
					Token:   common.BytesToAddress(event.Data[224:256]),
					PerUnit: big.NewInt(0).SetBytes(event.Data[256:288]),
				},
				GasGuarantee: pgas.GasPayment{
					Token:   common.BytesToAddress(event.Data[288:320]),
					PerUnit: big.NewInt(0).SetBytes(event.Data[320:352]),
				},
			}

			id, err := structs.WrapUint256(event.Topics[1][:])
			if err != nil {
				continue
			}

			PlanOrder(id, order)
		}
	}
}
