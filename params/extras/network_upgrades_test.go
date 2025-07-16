// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"testing"

	"github.com/luxfi/node/upgrade"
	"github.com/luxfi/node/upgrade/upgradetest"
	"github.com/luxfi/node/utils/constants"
	"github.com/luxfi/evm/utils"
	"github.com/stretchr/testify/require"
)

func TestNetworkUpgradesEqual(t *testing.T) {
	testcases := []struct {
		name      string
		upgrades1 *NetworkUpgrades
		upgrades2 *NetworkUpgrades
		expected  bool
	}{
		{
			name:      "nil upgrades",
			upgrades1: nil,
			upgrades2: nil,
			expected:  true,
		},
		{
			name:      "empty upgrades",
			upgrades1: &NetworkUpgrades{},
			upgrades2: &NetworkUpgrades{},
			expected:  true,
		},
		{
			name: "different subnet evm timestamp",
			upgrades1: &NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(1),
			},
			upgrades2: &NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(2),
			},
			expected: false,
		},
		{
			name: "same upgrades",
			upgrades1: &NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(1),
				DUpgradeTimestamp:  utils.NewUint64(2),
			},
			upgrades2: &NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(1),
				DUpgradeTimestamp:  utils.NewUint64(2),
			},
			expected: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.upgrades1.Equal(tt.upgrades2))
		})
	}
}

func TestCheckNetworkUpgradesCompatible(t *testing.T) {
	testcases := []struct {
		name      string
		upgrades  NetworkUpgrades
		time      uint64
		expected  error
	}{
		{
			name: "subnet evm used",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(100),
			},
			time:     101,
			expected: errUsedUpgrade,
		},
		{
			name: "subnet evm not used",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(100),
			},
			time:     99,
			expected: nil,
		},
		{
			name: "dUpgrade used",
			upgrades: NetworkUpgrades{
				DUpgradeTimestamp: utils.NewUint64(100),
			},
			time:     101,
			expected: errUsedUpgrade,
		},
		{
			name: "dUpgrade not used",
			upgrades: NetworkUpgrades{
				DUpgradeTimestamp: utils.NewUint64(100),
			},
			time:     99,
			expected: nil,
		},
		{
			name: "dUpgrade and subnet evm used",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(100),
				DUpgradeTimestamp:  utils.NewUint64(100),
			},
			time:     101,
			expected: errUsedUpgrade,
		},
		{
			name: "dUpgrade and subnet evm not used",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(100),
				DUpgradeTimestamp:  utils.NewUint64(100),
			},
			time:     99,
			expected: nil,
		},
		{
			name: "one used, one not used",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(100),
				DUpgradeTimestamp:  utils.NewUint64(200),
			},
			time:     101,
			expected: errUsedUpgrade,
		},
		{
			name: "cortinaX used",
			upgrades: NetworkUpgrades{
				CortinaXTime: utils.NewUint64(100),
			},
			time:     101,
			expected: errUsedUpgrade,
		},
		{
			name: "cortinaX not used",
			upgrades: NetworkUpgrades{
				CortinaXTime: utils.NewUint64(100),
			},
			time:     99,
			expected: nil,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.upgrades.CheckNetworkUpgradesCompatible(tt.time)
			require.ErrorIs(t, err, tt.expected)
		})
	}
}

func TestNetworkUpgradesValidation(t *testing.T) {
	testcases := []struct {
		name     string
		upgrades NetworkUpgrades
		expected error
	}{
		{
			name: "valid network upgrades",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(1),
				DUpgradeTimestamp:  utils.NewUint64(2),
			},
			expected: nil,
		},
		{
			name: "subnet evm nil",
			upgrades: NetworkUpgrades{
				DUpgradeTimestamp: utils.NewUint64(2),
			},
			expected: errCannotBeNil,
		},
		{
			name: "incompatible fork ordering",
			upgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(2),
				DUpgradeTimestamp:  utils.NewUint64(1),
			},
			expected: errIncompatibleForkSchedule,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.upgrades.Verify(constants.TestnetChainID, upgradetest.GetUpgradeDefaultTime("CortinaXTime"))
			require.ErrorIs(t, err, tt.expected)
		})
	}
}
