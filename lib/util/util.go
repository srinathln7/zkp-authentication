package util

import (
	"fmt"
	"math/big"
)

// ParseBigInt Parses a string and returns a pointer to the
// big.Int if successful
func ParseBigInt(str, param string) (*big.Int, error) {
	bigInt := new(big.Int)
	bigInt, valid := bigInt.SetString(str, 10)
	if !valid {
		return nil, fmt.Errorf("error parsing string %s to big.Int", param)
	}
	return bigInt, nil
}
