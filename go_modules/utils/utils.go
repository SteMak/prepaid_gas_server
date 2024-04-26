package utils

import (
	"math/big"
	"time"
)

var (
	err error
)

func UnixBig() *big.Int {
	return big.NewInt(time.Now().Unix())
}
