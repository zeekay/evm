// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"encoding/json"
	"math/big"
	"testing"
	"github.com/luxfi/evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestVerifyStateUpgrades(t *testing.T) {
	tests := []struct {
		name          string
		stateUpgrades []StateUpgrade
		expectedError string
	}{
		{
			name: "valid state upgrade",
			stateUpgrades: []StateUpgrade{
				{
					BlockNumber: big.NewInt(10),
					StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
						common.HexToAddress("0x1000000000000000000000000000000000000000"): {
							BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
						},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "invalid state upgrade - blackhole address",
			stateUpgrades: []StateUpgrade{
				{
					BlockNumber: big.NewInt(10),
					StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
						utils.BlackholeAddr: {
							BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
						},
					},
				},
			},
			expectedError: "cannot modify blackhole address",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Marshal and unmarshal to test JSON handling
			for i := range test.stateUpgrades {
				upgrade := &test.stateUpgrades[i]
				
				// Marshal to JSON
				data, err := json.Marshal(upgrade)
				require.NoError(t, err)

				// Unmarshal from JSON
				var unmarshaled StateUpgrade
				err = json.Unmarshal(data, &unmarshaled)
				require.NoError(t, err)

				// Verify the upgrade
				err = unmarshaled.verifyStateUpgrades(nil)
				if test.expectedError != "" {
					require.ErrorContains(t, err, test.expectedError)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestEqualStateUpgrade(t *testing.T) {
	addr1 := common.HexToAddress("0x1000000000000000000000000000000000000000")
	addr2 := common.HexToAddress("0x2000000000000000000000000000000000000000")

	tests := []struct {
		name     string
		upgrade1 *StateUpgrade
		upgrade2 *StateUpgrade
		expected bool
	}{
		{
			name:     "nil upgrades",
			upgrade1: nil,
			upgrade2: nil,
			expected: true,
		},
		{
			name:     "one nil upgrade",
			upgrade1: &StateUpgrade{},
			upgrade2: nil,
			expected: false,
		},
		{
			name: "equal upgrades",
			upgrade1: &StateUpgrade{
				BlockNumber: big.NewInt(10),
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					addr1: {
						BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
					},
				},
			},
			upgrade2: &StateUpgrade{
				BlockNumber: big.NewInt(10),
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					addr1: {
						BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
					},
				},
			},
			expected: true,
		},
		{
			name: "different addresses",
			upgrade1: &StateUpgrade{
				BlockNumber: big.NewInt(10),
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					addr1: {
						BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
					},
				},
			},
			upgrade2: &StateUpgrade{
				BlockNumber: big.NewInt(10),
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					addr2: {
						BalanceChange: (*math.HexOrDecimal256)(big.NewInt(1000)),
					},
				},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.upgrade1.Equal(test.upgrade2)
			require.Equal(t, test.expected, result)
		})
	}
}
