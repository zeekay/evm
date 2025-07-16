// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/luxfi/evm/precompile/modules"
	"github.com/luxfi/evm/precompile/precompileconfig"
	"github.com/luxfi/evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/luxfi/evm/precompile/contracts/deployerallowlist"
	"github.com/luxfi/evm/precompile/contracts/nativeminter"
	"github.com/luxfi/evm/precompile/contracts/txallowlist"
	"github.com/luxfi/evm/precompile/contracts/feemanager"
	"github.com/luxfi/evm/precompile/contracts/rewardmanager"
)

var errNoKey = errors.New("PrecompileUpgrade cannot be empty")
var errNoBlockTimestamp = errors.New("PrecompileUpgrade must have block timestamp")

// PrecompileUpgrade is a helper struct to define precompile upgrades.
type PrecompileUpgrade struct {
	precompileconfig.Config
	BlockTimestamp *uint64 `json:"blockTimestamp"`
	// Disable determines whether the precompile should be treated as disabled at the upgrade time.
	// If set to true, this overrides the timestamp set in the Config.
	// A disabled precompile will be removed from the list of active precompiles configured by the
	// node at the upgrade time. If a fork scheduled after the upgrade re-enables the precompile,
	// it will be restored to the list of active precompiles.
	Disable bool `json:"disable,omitempty"`
	// a single precompile is chosen by setting exactly one of these
	// Deprecated: use Config.Key() instead.
	PrecompileIdentifier
}

type PrecompileIdentifier struct {
	ContractDeployerAllowListConfig *deployerallowlist.Config `json:"contractDeployerAllowListConfig,omitempty"`
	ContractNativeMinterConfig      *nativeminter.Config      `json:"contractNativeMinterConfig,omitempty"`
	TxAllowListConfig               *txallowlist.Config       `json:"txAllowListConfig,omitempty"`
	FeeManagerConfig                *feemanager.Config        `json:"feeManagerConfig,omitempty"`
	RewardManagerConfig             *rewardmanager.Config     `json:"rewardManagerConfig,omitempty"`
	// ADD YOUR PRECOMPILE HERE
}

// UnmarshalJSON unmarshals the json into the correct precompile config type
// Ex: {"type": "txAllowListConfig", "adminAddresses": [address1, address2, ...]}
func (p *PrecompileUpgrade) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check if the upgrade is marked as disabled
	if disableData, ok := raw["disable"]; ok {
		if err := json.Unmarshal(disableData, &p.Disable); err != nil {
			return err
		}
	}

	// Check for blockTimestamp
	if blockTimestampData, ok := raw["blockTimestamp"]; ok {
		var blockTimestamp uint64
		if err := json.Unmarshal(blockTimestampData, &blockTimestamp); err != nil {
			return err
		}
		p.BlockTimestamp = &blockTimestamp
	}

	module, ok := raw["type"]
	if !ok {
		return errors.New("missing 'type' field")
	}

	var typeStr string
	if err := json.Unmarshal(module, &typeStr); err != nil {
		return err
	}

	delete(raw, "type")
	delete(raw, "disable")
	delete(raw, "blockTimestamp")

	cfgData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	key := modules.ReservedAddress(typeStr)

	cfg := modules.GetPrecompileModuleByAddress(key)
	if cfg == nil {
		return fmt.Errorf("unknown type: %s", typeStr)
	}

	p.Config = cfg
	if err := json.Unmarshal(cfgData, p.Config); err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals the precompile upgrade into json
func (p *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	if p.Config != nil {
		configData := p.Config

		// Marshal the config to JSON and then unmarshal into map
		jsonData, err := json.Marshal(configData)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(jsonData, &data); err != nil {
			return nil, err
		}

		// Add the type field
		data["type"] = modules.GetPrecompileModuleByAddress(p.Config.Key()).Address
	}

	if p.Disable {
		data["disable"] = true
	}

	if p.BlockTimestamp != nil {
		data["blockTimestamp"] = *p.BlockTimestamp
	}

	return json.Marshal(data)
}

// VerifyPrecompileUpgrade checks that [pu] is a valid upgrade.
func (pu *PrecompileUpgrade) VerifyPrecompileUpgrade(previous *PrecompileUpgrade, current *UpgradeableRules) error {
	if pu.BlockTimestamp == nil {
		return errNoBlockTimestamp
	}

	if pu.Disable && pu.Config == nil {
		return fmt.Errorf("PrecompileUpgrade cannot disable an empty upgrade")
	}

	// Verify specified config
	switch {
	case pu.Config != nil:
		if err := pu.Config.Verify(previous.GetPrecompileConfig(pu.Config.Key(), current), current); err != nil {
			return err
		}
	default:
		return errNoKey
	}

	// If this upgrade disables the precompile, verify that the precompile is active at the timestamp
	if pu.Disable && !pu.Config.IsDisabled() {
		return fmt.Errorf("PrecompileUpgrade cannot disable a precompile with a mismatched timestamp: the precompile must be configured at the same timestamp as the upgrade")
	}

	// Verify the upgrade does not configure a precompile that is already configured.
	if previous.GetPrecompileConfig(pu.Config.Key(), current) != nil && !previous.GetPrecompileConfig(pu.Config.Key(), current).IsDisabled() && !pu.Disable {
		return fmt.Errorf("PrecompileUpgrade cannot configure a precompile that is already configured")
	}

	return nil
}

// GetPrecompileConfig returns the precompile config for the given key
func (pu *PrecompileUpgrade) GetPrecompileConfig(key common.Address, chainConfig *UpgradeableRules) precompileconfig.Config {
	if pu.Config != nil && pu.Config.Key() == key {
		return pu.Config
	}
	return nil
}

// Equal returns true if the PrecompileUpgrade is equal to the other PrecompileUpgrade.
func (pu *PrecompileUpgrade) Equal(other *PrecompileUpgrade) bool {
	if pu == other {
		return true
	}
	if pu == nil || other == nil {
		return false
	}
	if pu.Disable != other.Disable {
		return false
	}
	if (pu.BlockTimestamp == nil) != (other.BlockTimestamp == nil) {
		return false
	}
	if pu.BlockTimestamp != nil && *pu.BlockTimestamp != *other.BlockTimestamp {
		return false
	}
	if (pu.Config == nil) != (other.Config == nil) {
		return false
	}
	if pu.Config != nil && !pu.Config.Equal(other.Config) {
		return false
	}
	return true
}
