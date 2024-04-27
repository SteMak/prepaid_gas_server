package executor

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

func monitorMessages() {
	for {
		result, _ := db.GetMessages(false, offset, 1000)
		offset += uint64(len(result))

		for _, item := range result {
			message, sign, _ := structs.UnwrapDBMessage(item)
			go planMessage(message, sign)
		}

		time.Sleep(time.Second * 5)
	}
}

func planMessage(message structs.Message, sign structs.Signature) {
	order := orders[hex.EncodeToString(message.Order[:])]
	if order == nil {
		return
	}

	// start + window < now
	if big.NewInt(0).Add(message.Start.ToBig(), order.TxWindow).Cmp(utils.UnixBig()) == -1 {
		return
	}

	// start - now > delay
	if big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()).Cmp(big.NewInt(int64(delay))) == 1 {
		time.Sleep(time.Second * time.Duration(
			big.NewInt(0).Sub(
				big.NewInt(0).Sub(
					message.Start.ToBig(), utils.UnixBig(),
				),
				big.NewInt(int64(delay)),
			).Int64(),
		))
		used, _ := onchain.PGas.Nonce(nil, common.Address(message.From), message.Nonce.ToBig())
		if used {
			return
		}
	}

	time.Sleep(time.Second * time.Duration(big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()).Int64()))

	_, err = onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
	for err != nil && big.NewInt(0).Add(message.Start.ToBig(), order.TxWindow).Cmp(utils.UnixBig()) == 1 {
		time.Sleep(time.Second)
		_, err = onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
	}
}
