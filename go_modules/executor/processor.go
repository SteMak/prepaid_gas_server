package executor

import (
	"math/big"
	"time"

	"github.com/prepaidGas/prepaidgas-server/go_modules/db"
	"github.com/prepaidGas/prepaidgas-server/go_modules/structs"
)

func Processor() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		result, _ := db.GetMessages(false, last_message, 1000)
		last_message += uint64(len(result))

		for _, item := range result {
			if orders[big.NewInt(0).SetBytes(item.Order[:])] != nil {
				message, sign, _ := structs.UnwrapDBMessage(item)
				go PlanMessage(message, sign)
			}
		}
	}
}
