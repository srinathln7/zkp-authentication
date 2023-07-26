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

// StringToUniqueBigInt: Converts the given string uniquely to a Big Int
// assuming the given string consists of only ASCII character set
func StringToUniqueBigInt(input string) *big.Int {
	base := big.NewInt(256) // 256 is used as the base, assuming ASCII character set (8-bit)

	var result big.Int
	for _, ch := range input {
		result.Mul(&result, base)                  // Shift left by base
		result.Add(&result, big.NewInt(int64(ch))) // Add the character value to the result
	}

	return &result
}
