// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"
	"github.com/luxfi/evm/params"
	"github.com/luxfi/evm/params/extras"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/triedb"
)

// Configure applies the state upgrade to the state.
func Configure(stateUpgrade *extras.StateUpgrade, chainConfig ChainContext, state StateDB, blockContext BlockContext) error {
	isEIP158 := chainConfig.IsEIP158(blockContext.Number())
	for account, upgrade := range stateUpgrade.StateUpgradeAccounts {
		if err := upgradeAccount(account, upgrade, state, isEIP158); err != nil {
			return err
		}
	}
	return nil
}

// upgradeAccount applies the modifications in [upgrade] to [account].
func upgradeAccount(account common.Address, upgrade extras.StateUpgradeAccount, state StateDB, isEIP158 bool) error {
	// Change the account balance
	if upgrade.BalanceChange != nil {
		state.AddBalance(account, upgrade.BalanceChange.ToInt(), triedb.BalanceChangeUnspecified)
	}

	// Change the code
	if upgrade.Code != nil {
		state.SetCode(account, *upgrade.Code)
	}

	// Update storage
	for key, value := range upgrade.Storage {
		state.SetState(account, key, value)
	}

	// Create the account if it does not exist
	if !state.Exist(account) && !isEIP158 {
		state.CreateAccount(account)
	}

	return nil
}

// ChainContext defines an interface that provides information about the chain configuration.
type ChainContext interface {
	IsEIP158(num *big.Int) bool
}

// StateDB defines an interface for interacting with the EVM state during upgrades.
type StateDB interface {
	AddBalance(addr common.Address, amount *big.Int, reason triedb.BalanceChangeReason)
	SetCode(addr common.Address, code []byte)
	SetState(addr common.Address, key common.Hash, value common.Hash)
	Exist(addr common.Address) bool
	CreateAccount(addr common.Address)
}

// BlockContext defines an interface that provides information about the block being processed.
type BlockContext interface {
	Number() *big.Int
	Timestamp() uint64
}
