package tps

import (
	"context"
)

type Client interface {
	LatestBlockHeight(context.Context) (uint64, error)
	CountTx(context.Context, uint64) (int, error)
	CountPendingTx(context.Context) (int, error)
	Nonce(context.Context, string) (uint64, error)
}
