package eth

import (
	"context"
	"crypto/ecdsa"
)

type ContextKey int

const (
	PrivKeyCtxKey ContextKey = iota
)

func CtxWithPriv(ctx context.Context, priv *ecdsa.PrivateKey) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, PrivKeyCtxKey, priv)
}

func GetPriv(ctx context.Context) *ecdsa.PrivateKey {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	priv, ok := ctx.Value(PrivKeyCtxKey).(*ecdsa.PrivateKey)
	if !ok {
		return nil
	}

	return priv
}
