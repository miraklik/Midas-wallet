package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func CheckAddress(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}

	return nil
}

func FormatAmount(amount *big.Int) string {
	return fmt.Sprintf("%s WEI", amount.String())
}
