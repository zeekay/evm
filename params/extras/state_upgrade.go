// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"fmt"
	"math/big"
	"reflect"
	"github.com/luxfi/evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
)

// StateUpgrade describes the modifications to be made to the state during
// a hard fork.
type StateUpgrade struct {
	BlockNumber *big.Int `json:"blockNumber,omitempty"`

	// map from account address to the modification to be made to the account.
	StateUpgradeAccounts map[common.Address]StateUpgradeAccount `json:"accounts"`
}

// StateUpgradeAccount describes the modifications to be made to an account during
// a hard fork.
type StateUpgradeAccount struct {
	Code         *hexutil.Bytes                        `json:"code,omitempty"`
	Storage      map[common.Hash]common.Hash           `json:"storage,omitempty"`
	BalanceChange *math.HexOrDecimal256                 `json:"balanceChange,omitempty"`
}

func (s *StateUpgrade) Equal(other *StateUpgrade) bool {
	return reflect.DeepEqual(s, other)
}

// verifyStateUpgrades checks that the configuration is valid.
func (s *StateUpgrade) verifyStateUpgrades(config *ChainConfig) error {
	if s == nil {
		return nil
	}

	for address := range s.StateUpgradeAccounts {
		if err := utils.VerifyAddress(address); err != nil {
			return fmt.Errorf("invalid address: %s", address)
		}
	}

	return nil
}
