// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"encoding/json"
	"math/big"
	"testing"
	"github.com/luxfi/evm/commontype"
	"github.com/luxfi/evm/precompile/contracts/deployerallowlist"
	"github.com/luxfi/evm/precompile/contracts/feemanager"
	"github.com/luxfi/evm/precompile/contracts/nativeminter"
	"github.com/luxfi/evm/precompile/contracts/rewardmanager"
	"github.com/luxfi/evm/precompile/contracts/txallowlist"
	"github.com/luxfi/evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyWithChainConfig(t *testing.T) {
	admins := []common.Address{{1}}
	baseConfig := ChainConfig{
		UpgradeableRules: UpgradeableRules{
			HomesteadBlock: big.NewInt(0),
		},
	}
	config := &baseConfig
	config.PrecompileUpgrade = PrecompileUpgrade{
		TxAllowListConfig: txallowlist.NewConfig(utils.NewUint64(10), admins, nil, nil),
	}

	// Check this config is valid
	err := config.Verify()
	require.NoError(t, err)

	// Check the precompile is configured correctly
	require.True(t, config.IsPrecompileEnabled(txallowlist.ContractAddress, big.NewInt(20)))
	require.False(t, config.IsPrecompileEnabled(txallowlist.ContractAddress, big.NewInt(5)))

	// Check that the upgrade is not valid when applied at the same timestamp
	upgrade := PrecompileUpgrade{
		Config: txallowlist.NewConfig(utils.NewUint64(10), admins, nil, nil),
	}
	err = upgrade.VerifyPrecompileUpgrade(&config.PrecompileUpgrade, &config.UpgradeableRules)
	require.Error(t, err)

	// Check that the upgrade is valid when applied at a later timestamp
	upgrade.Config = txallowlist.NewConfig(utils.NewUint64(20), admins, nil, nil)
	err = upgrade.VerifyPrecompileUpgrade(&config.PrecompileUpgrade, &config.UpgradeableRules)
	require.NoError(t, err)
}

func TestEqualPrecompileUpgrade(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	managers := []common.Address{{3}}

	tests := []struct {
		name     string
		upgrade1 PrecompileUpgrade
		upgrade2 PrecompileUpgrade
		expected bool
	}{
		{
			name:     "empty upgrades",
			upgrade1: PrecompileUpgrade{},
			upgrade2: PrecompileUpgrade{},
			expected: true,
		},
		{
			name:     "tx allowlist",
			upgrade1: PrecompileUpgrade{Config: txallowlist.NewConfig(utils.NewUint64(1), admins, enableds, managers)},
			upgrade2: PrecompileUpgrade{Config: txallowlist.NewConfig(utils.NewUint64(1), admins, enableds, managers)},
			expected: true,
		},
		{
			name:     "different precompiles",
			upgrade1: PrecompileUpgrade{Config: txallowlist.NewConfig(utils.NewUint64(1), admins, nil, nil)},
			upgrade2: PrecompileUpgrade{Config: deployerallowlist.NewConfig(utils.NewUint64(1), admins, nil, nil)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.upgrade1.Equal(&tt.upgrade2))
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	admins := []common.Address{{1}}
	feeConfig := commontype.FeeConfig{
		GasLimit:    big.NewInt(8_000_000),
		TargetBlockRate: 2,
		MinBaseFee:      big.NewInt(25_000_000_000),
		TargetGas:       big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),
		MinBlockGasCost:          big.NewInt(0),
		MaxBlockGasCost:          big.NewInt(1_000_000),
		BlockGasCostStep:         big.NewInt(200_000),
	}

	tests := []struct {
		name        string
		upgrade     PrecompileUpgrade
		expectedJSON string
	}{
		{
			name:         "tx allowlist",
			upgrade:      PrecompileUpgrade{Config: txallowlist.NewConfig(utils.NewUint64(1), admins, nil, nil)},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000002","adminAddresses":["0x0100000000000000000000000000000000000000"],"timestamp":1}`,
		},
		{
			name:         "fee manager",
			upgrade:      PrecompileUpgrade{Config: feemanager.NewConfig(utils.NewUint64(1), admins, nil, &feeConfig)},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000003","adminAddresses":["0x0100000000000000000000000000000000000000"],"initialFeeConfig":{"gasLimit":8000000,"targetBlockRate":2,"minBaseFee":25000000000,"targetGas":15000000,"baseFeeChangeDenominator":36,"minBlockGasCost":0,"maxBlockGasCost":1000000,"blockGasCostStep":200000},"timestamp":1}`,
		},
		{
			name:         "deployer allowlist",
			upgrade:      PrecompileUpgrade{Config: deployerallowlist.NewConfig(utils.NewUint64(1), admins, nil, nil)},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000001","adminAddresses":["0x0100000000000000000000000000000000000000"],"timestamp":1}`,
		},
		{
			name:         "native minter",
			upgrade:      PrecompileUpgrade{Config: nativeminter.NewConfig(utils.NewUint64(1), admins, nil, nil, nil)},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000004","adminAddresses":["0x0100000000000000000000000000000000000000"],"timestamp":1}`,
		},
		{
			name:         "reward manager",
			upgrade:      PrecompileUpgrade{Config: rewardmanager.NewConfig(utils.NewUint64(1), admins, nil, nil, nil)},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000005","adminAddresses":["0x0100000000000000000000000000000000000000"],"timestamp":1}`,
		},
		{
			name:         "tx allowlist disabled",
			upgrade:      PrecompileUpgrade{Config: txallowlist.NewConfig(utils.NewUint64(1), admins, nil, nil), Disable: true},
			expectedJSON: `{"type":"0x0200000000000000000000000000000000000002","adminAddresses":["0x0100000000000000000000000000000000000000"],"disable":true,"timestamp":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.upgrade)
			require.NoError(t, err)
			require.JSONEq(t, tt.expectedJSON, string(data))

			var unmarshaled PrecompileUpgrade
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)
			require.Equal(t, tt.upgrade, unmarshaled)
		})
	}
}
