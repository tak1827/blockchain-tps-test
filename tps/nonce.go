package tps

import (
	"context"
	"sync/atomic"
)

type Nonce struct {
	current uint64
}

func NewNonce(ctx context.Context, client Client, addr string) (Nonce, error) {
	current, err := client.Nonce(ctx, addr)
	if err != nil {
		return Nonce{}, err
	}

	return Nonce{current: current}, nil
}

func (n *Nonce) Increment() uint64 {
	return atomic.AddUint64(&n.current, 1) - 1
}

func (n *Nonce) Reset(nonce uint64) {
	atomic.StoreUint64(&n.current, nonce)
}

func (n *Nonce) Current() uint64 {
	return atomic.LoadUint64(&n.current)
}
