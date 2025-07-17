// (c) 2025, Lux Industries, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import "time"


// TimePtrToNewUint64 converts a pointer to a time to a pointer to a uint64
func TimePtrToNewUint64(timePtr *time.Time) *uint64 {
	if timePtr == nil {
		return nil
	}
	return TimeToNewUint64(*timePtr)
}