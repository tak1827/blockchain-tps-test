package tps

import (
	"sync/atomic"
	"time"
)

func NextIdlingDuration(idlingDuration *uint32, txs, pendingTxs uint32) {
	current := atomic.LoadUint32(idlingDuration)
	next := uint32(0)
	if pendingTxs/txs > 1 {
		next = (current + (pendingTxs / txs)) * uint32(1*time.Millisecond)
	}
	atomic.StoreUint32(idlingDuration, next)
}

func ToDuration(idlingDuration *uint32) time.Duration {
	return time.Duration(int64(atomic.LoadUint32(idlingDuration)))
}
