// (c) 2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// VerifyAddress checks if the given string is a valid Ethereum address
func VerifyAddress(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid hex address: %s", address)
	}
	return nil
}