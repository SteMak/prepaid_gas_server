package executor

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/onchain"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
	"github.com/prepaidGas/prepaidgas-server/go_modules/utils"
)

var (
	delay *big.Int
)

func initMonitor(prevalidate_delay uint32) {
	delay = big.NewInt(int64(prevalidate_delay))
}

func monitor() {
	for {
		result, err := db.GetMessages(false, offset, 1000)
		if err != nil {
			log.Printf("monitor db: %s\n\n", err.Error())
			time.Sleep(time.Second)
			continue
		}
		offset += uint64(len(result))

		for _, item := range result {
			message, sign, _ := structs.UnwrapDBMessage(item)
			go planMessage(message, sign)
		}

		time.Sleep(time.Second * 5)
	}
}

func planMessage(message structs.Message, sign structs.Signature) {
	order := orders[message.Order.ToString()]
	if order == nil {
		log.Printf("message order not promised: \"%#v\"\n\n", message)
		return
	}

	// start + window < now
	if big.NewInt(0).Add(message.Start.ToBig(), order.TxWindow).Cmp(utils.UnixBig()) == -1 {
		log.Printf("message in past: \"%#v\"\n\n", message)
		return
	}

	// start - now > delay
	if big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()).Cmp(delay) == 1 {
		sleep := big.NewInt(0).Sub(big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig()), delay)
		if !sleep.IsInt64() {
			log.Printf("message sleep time out of life: \"%#v\" %s\n\n", message, sleep.Text(16))
			return
		}

		log.Printf("message nonce check planned: \"%#v\"\n\n", message)
		time.Sleep(time.Second * time.Duration(sleep.Int64()))

		used, _ := onchain.PGas.Nonce(nil, common.Address(message.From), message.Nonce.ToBig())
		if used {
			log.Printf("message nonce already used: \"%#v\"\n\n", message)
			return
		}
	}

	sleep := big.NewInt(0).Sub(message.Start.ToBig(), utils.UnixBig())
	if !sleep.IsInt64() {
		log.Printf("message sleep time out of life: \"%#v\" %s\n\n", message, sleep.Text(16))
		return
	}

	log.Printf("message planned: \"%#v\"\n\n", message)
	time.Sleep(time.Second * time.Duration(sleep.Int64()))

	_, err := onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
	for err != nil {
		log.Printf("message execute: \"%#v\" %s\n\n", message, err.Error())

		time.Sleep(time.Second)
		if big.NewInt(0).Add(message.Start.ToBig(), order.TxWindow).Cmp(utils.UnixBig()) == -1 {
			return
		}

		_, err = onchain.PGas.Execute(onchain.Transactor, onchain.WrapPGasMessage(message), sign.ToOnchain())
	}

	log.Printf("message success execute: \"%#v\"\n\n", message)
}
