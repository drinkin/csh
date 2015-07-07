package csh

import "time"

// Helper to check for optional expire time
func hasExpire(exps []time.Duration) bool {
	return len(exps) > 0 && exps[0] > 0
}
