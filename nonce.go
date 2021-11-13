package main

import (
	"context"
	"sync/atomic"

	"github.com/min-sys/tracking-contract/cli/eth"
)

type Nonce struct {
	current uint64
}

func NewNonce(ctx context.Context, client *eth.Client) (Nonce, error) {
	current, err := client.Nonce(ctx)
	if err != nil {
		return Nonce{}, err
	}

	return Nonce{current: current}, nil
}

func (n *Nonce) Increment() uint64 {
	current := atomic.LoadUint64(&n.current)
	atomic.AddUint64(&n.current, 1)
	return current
}

func (n *Nonce) Reset(nonce uint64) {
	atomic.StoreUint64(&n.current, nonce)
}

func (n *Nonce) Current() uint64 {
	return atomic.LoadUint64(&n.current)
}
