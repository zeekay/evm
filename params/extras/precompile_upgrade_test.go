// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"math/big"
	"testing"
	"github.com/luxfi/evm/precompile/contracts/deployerallowlist"
	"github.com/luxfi/evm/precompile/contracts/txallowlist"
	"github.com/luxfi/evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyUpgradeConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	managers := []common.Address{{3}}
	chainConfig := getChainConfig(admins, enableds, managers)

	type test struct {
		upgrades []PrecompileUpgrade
		expected error
	}

	var tests = []test{
		{
			upgrades: []PrecompileUpgrade{
				{BlockTimestamp: utils.NewUint64(1)},
			},
			expected: errNoKey,
		},
		{
			upgrades: []PrecompileUpgrade{
				{
					Config: deployerallowlist.NewConfig(utils.NewUint64(3), admins, enableds, managers),
				},
			},
			expected: errNoBlockTimestamp,
		},
		{
			upgrades: []PrecompileUpgrade{
				{
					BlockTimestamp: utils.NewUint64(5),
					Config:         txallowlist.NewConfig(utils.NewUint64(10), admins, enableds, managers),
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		for i := range tt.upgrades {
			err := tt.upgrades[i].VerifyPrecompileUpgrade(&chainConfig.PrecompileUpgrade, &chainConfig.UpgradeableRules)
			require.ErrorIs(t, err, tt.expected)
		}
	}
}

func getChainConfig(admins []common.Address, enableds []common.Address, managers []common.Address) ChainConfig {
	return ChainConfig{
		UpgradeableRules: UpgradeableRules{
			HomesteadBlock: big.NewInt(0),
		},
		PrecompileUpgrade: PrecompileUpgrade{
			TxAllowListConfig: txallowlist.NewConfig(utils.NewUint64(10), admins, enableds, managers),
		},
	}
}
