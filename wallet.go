package main

import (
	"crypto/ecdsa"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/min-sys/tracking-contract/cli/eth"
	"github.com/pkg/errors"
)

type Wallet struct {
	privs      []*ecdsa.PrivateKey
	nonces     map[*ecdsa.PrivateKey]*Nonce
	roteteSlot uint32
}

func NewWallet(client *eth.Client, keys []string) (Wallet, error) {
	var (
		privs  = make([]*ecdsa.PrivateKey, len(keys))
		nonces = make(map[*ecdsa.PrivateKey]*Nonce, len(keys))
		err    error
	)

	for i := range keys {
		if privs[i], err = crypto.HexToECDSA(keys[i]); err != nil {
			return Wallet{}, errors.Wrap(err, "err HexToECDSA")
		}

		ctx := eth.CtxWithPriv(nil, privs[i])
		nonce, err := NewNonce(ctx, client)
		if err != nil {
			return Wallet{}, errors.Wrap(err, "err NewNonce")
		}
		nonces[privs[i]] = &nonce
	}

	return Wallet{
		privs:  privs,
		nonces: nonces,
	}, nil
}

func (w *Wallet) RotatePriv() *ecdsa.PrivateKey {
	slot := atomic.LoadUint32(&w.roteteSlot)
	atomic.AddUint32(&w.roteteSlot, 1)
	return w.privs[slot%uint32(len(w.privs))]
}

func (w *Wallet) IncrementNonce(priv *ecdsa.PrivateKey) uint64 {
	return w.nonces[priv].Increment()
}

func (w *Wallet) CurrentNonce(priv *ecdsa.PrivateKey) uint64 {
	return w.nonces[priv].Current()
}

func (w *Wallet) RecetNonce(priv *ecdsa.PrivateKey, nonce uint64) {
	w.nonces[priv].Reset(nonce)
}
