package utils

import (
	"math/big"

	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain/pgas"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func IsOrderRisky(id structs.Uint256, order pgas.Order) bool {
	messages, err := db.GetMessagesByOrder(id, 0, 1)
	if err != nil || uint64(len(messages)) > 0 {
		return true
	}

	if order.GasGuarantee.PerUnit.Cmp(big.NewInt(0)) == 0 {
		return false
	}

	if order.GasGuarantee.PerUnit.Cmp(big.NewInt(100001)) == -1 &&
		order.Gas.Cmp(big.NewInt(1000001)) == -1 &&
		big.NewInt(0).Add(UnixBig(), big.NewInt(60*60*30)).Cmp(order.End) == 1 {
		return false
	}

	return true
}
