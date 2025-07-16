// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"encoding/json"

	"github.com/luxfi/evm/precompile/modules"
	"github.com/luxfi/evm/precompile/precompileconfig"
)

type Precompiles map[string]precompileconfig.Config

// UnmarshalJSON parses the JSON-encoded data into the ChainConfigPrecompiles.
// ChainConfigPrecompiles is a map of precompile module keys to their
// configuration.
func (ccp *Precompiles) UnmarshalJSON(data []byte) error {
	*ccp = make(Precompiles)

	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for module, config := range raw {
		m, ok := modules.GetPrecompileModuleByKey(module)
		if !ok {
			return modules.ErrUndefinedPrecompileKey{
				Key: module,
			}
		}

		moduleCfg := m.MakeConfig()
		if err := json.Unmarshal(config, moduleCfg); err != nil {
			return err
		}

		(*ccp)[module] = moduleCfg
	}

	return nil
}
