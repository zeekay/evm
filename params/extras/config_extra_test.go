// (c) 2024 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"testing"

	"github.com/luxfi/evm/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsTimestampForked(t *testing.T) {
	type test struct {
		fork     *uint64
		block    uint64
		isForked bool
	}

	zero := uint64(0)
	five := uint64(5)
	tests := []test{
		{fork: nil, block: 0, isForked: false},
		{fork: &zero, block: 0, isForked: true},
		{fork: &zero, block: 1, isForked: true},
		{fork: &five, block: 4, isForked: false},
		{fork: &five, block: 5, isForked: true},
		{fork: &five, block: 6, isForked: true},
	}

	for i, tt := range tests {
		isForked := utils.IsTimestampForked(tt.fork, tt.block)
		assert.Equal(t, tt.isForked, isForked, "test %d", i)
	}
}
