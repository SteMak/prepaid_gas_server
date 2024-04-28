package utils

import (
	"math/big"
	"time"
)

func UnixBig() *big.Int {
	return big.NewInt(time.Now().Unix())
}
