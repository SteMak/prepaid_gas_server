package executor

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func Acceptor(pgas_address structs.Address) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.BytesToAddress(pgas_address[:])},
		Topics:    [][]common.Hash{{common.BytesToHash(crypto.Keccak256([]byte("OrderCreate(uint256,uint256)")))}},
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
			fmt.Printf("%+v\n", event)
		}
	}
}
