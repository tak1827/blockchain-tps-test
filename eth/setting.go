package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"
)

type Setting struct {
	Endpoint string
	Timeout  time.Duration

	PrivKey  *ecdsa.PrivateKey
	GasPrice *big.Int
}

func (s Setting) TimeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, s.Timeout)
}
